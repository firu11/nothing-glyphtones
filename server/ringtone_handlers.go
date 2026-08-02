package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"glyphtones/database"
	"glyphtones/templates/components"
	"glyphtones/templates/views"
	"glyphtones/utils"

	"github.com/labstack/echo/v5"
)

const downloadedCookieTTL = 1_000_000 * time.Hour

func (s *Server) renameView(c *echo.Context) error {
	authorID := s.currentUserID(c)
	if authorID == 0 {
		return Render(c, views.OtherErrorView(http.StatusBadRequest, errors.New("You're not logged in.")))
	}

	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	ringtone, err := s.store.GetRingtone(c.Request().Context(), displayID, authorID)
	if err != nil {
		return internalError(c, "get ringtone for rename", err, false)
	}

	return Render(c, components.Rename(ringtone, nil))
}

func (s *Server) rename(c *echo.Context) error {
	authorID := s.currentUserID(c)
	if authorID == 0 {
		return errors.New("You're not logged in.")
	}
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	newName := c.FormValue("name")
	if !ringtoneNameR.MatchString(newName) {
		ringtone, err := s.store.GetRingtone(c.Request().Context(), displayID, authorID)
		if err != nil {
			return internalError(c, "get ringtone after invalid rename", err, true)
		}
		ringtone.Name = newName
		return Render(c, components.Rename(ringtone, errors.New("The name must be 2-20 letters and only a-z and some special characters.")))
	}
	if err := s.store.RenameRingtone(c.Request().Context(), displayID, newName, authorID); err != nil {
		slog.Error("failed to rename ringtone", "method", c.Request().Method, "path", c.Request().URL.Path, "display_id", displayID, "error", err)
		ringtone, getErr := s.store.GetRingtone(c.Request().Context(), displayID, authorID)
		if getErr != nil {
			return internalError(c, "reload ringtone after failed rename", getErr, true)
		}
		return Render(c, components.Rename(ringtone, errors.New("Something went wrong")))
	}
	ringtone, err := s.store.GetRingtone(c.Request().Context(), displayID, authorID)
	if err != nil {
		return internalError(c, "reload renamed ringtone", err, true)
	}

	var effects []database.EffectModel
	if authorID == 1 {
		effects, err = s.store.GetEffects(c.Request().Context())
		if err != nil {
			return internalError(c, "query effects after rename", err, false)
		}
	}
	return Render(c, components.Captions(ringtone, effects, true, authorID == 1))
}

func (s *Server) updateRingtoneMetadata(c *echo.Context) error {
	if s.currentUserID(c) != 1 {
		return c.NoContent(http.StatusForbidden)
	}

	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	category, categoryErr := strconv.Atoi(c.FormValue("category"))
	effectID, effectErr := strconv.Atoi(c.FormValue("effect"))
	if categoryErr != nil || category < 1 || category > len(components.Categories) || effectErr != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	effects, err := s.store.GetEffects(ctx)
	if err != nil {
		return internalError(c, "query effects for metadata update", err, false)
	}
	validEffect := false
	for _, effect := range effects {
		if effect.ID == effectID {
			validEffect = true
			break
		}
	}
	if !validEffect {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.store.UpdateRingtoneMetadata(ctx, displayID, category, effectID, c.FormValue("auto_generated") == "on"); err != nil {
		return internalError(c, "update ringtone metadata", err, false)
	}
	ringtone, err := s.store.GetRingtone(ctx, displayID, 1)
	if err != nil {
		return internalError(c, "reload updated ringtone", err, false)
	}
	return Render(c, components.Captions(ringtone, effects, true, true))
}

func (s *Server) uploadView(c *echo.Context) error {
	effects, err := s.store.GetEffects(c.Request().Context())
	if err != nil {
		return internalError(c, "query upload effects", err, true)
	}

	return Render(c, views.Upload(s.loggedInFromCookie(c), c.FormValue("c"), effects, "", "", nil))
}

func (s *Server) uploadFile(c *echo.Context) error {
	authorID := s.currentUserID(c)
	if authorID == 0 {
		return Render(c, views.OtherError(http.StatusBadRequest, errors.New("Only logged-in authors can upload Glyphtones")))
	}
	ctx := c.Request().Context()
	author, err := s.store.GetAuthor(ctx, authorID)
	if err != nil {
		return internalError(c, "get uploading author", err, false)
	}

	errorHandler := func(mainErr error) error {
		effects, err := s.store.GetEffects(ctx)
		if err != nil {
			return internalError(c, "query effects for upload error", err, true)
		}
		return Render(c, views.UploadForm(c.FormValue("c"), effects, c.FormValue("e"), c.FormValue("name"), true, mainErr))
	}

	if author.Banned {
		return errorHandler(errors.New("You cannot upload since you are banned!"))
	}

	name := c.FormValue("name")
	if !ringtoneNameR.MatchString(name) {
		return errorHandler(errors.New("Name must be 2-30 characters long and without diacritics."))
	}
	category, err1 := strconv.Atoi(c.FormValue("c"))
	effect, err2 := strconv.Atoi(c.FormValue("e"))
	if err1 != nil || err2 != nil || c.FormValue("gen") == "" {
		return errorHandler(errors.New("Missing form values."))
	}
	autoGenerated := c.FormValue("gen") == "true"
	file, err := c.FormFile("ringtone")
	if err != nil {
		return errorHandler(errors.New("Missing the file."))
	}
	if !hasOGGExtension(file.Filename) {
		return errorHandler(errors.New("It seems that the file provided is not a Nothing Glyphtone."))
	}

	req := uploadRequest{
		Name:          name,
		Category:      category,
		Effect:        effect,
		AuthorID:      authorID,
		AutoGenerated: autoGenerated,
	}
	if err := s.saveUploadedRingtone(ctx, req, file); err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidRingtoneFile):
			return errorHandler(errors.New("It seems that the file provided is not a Nothing Glyphtone."))
		case errors.Is(err, errDuplicateRingtone):
			return Render(c, views.OtherError(http.StatusBadRequest, errors.New("You're trying to upload a file which has been uploaded before. Please do not do that...")))
		case errors.Is(err, errFileTooLarge):
			return Render(c, views.OtherError(http.StatusBadRequest, err))
		default:
			return internalError(c, "save uploaded ringtone", err, false)
		}
	}

	return Render(c, views.SuccessfulUpload())
}

func (s *Server) ringtoneGlyphs(c *echo.Context) error {
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}

	glyphs, err := s.store.GetRingtoneGlyphs(c.Request().Context(), displayID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		return internalError(c, "get ringtone glyphs", err, false)
	}
	if len(glyphs) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000, immutable")
	return c.Blob(http.StatusOK, "application/octet-stream", glyphs)
}

func (s *Server) downloadRingtone(c *echo.Context) error {
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	_, err := c.Cookie(fmt.Sprintf("Glyphtone_%s_downloaded", displayID))
	if err == nil {
		return c.NoContent(http.StatusOK)
	}
	if err := s.store.RingtoneIncreaseDownload(c.Request().Context(), displayID); err != nil {
		return internalError(c, "increment ringtone download", err, false)
	}
	cookie := http.Cookie{
		Name:     fmt.Sprintf("Glyphtone_%s_downloaded", displayID),
		Value:    "true",
		Expires:  time.Now().Add(downloadedCookieTTL),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.Production,
	}
	c.SetCookie(&cookie)
	return c.NoContent(http.StatusOK)
}

func (s *Server) deleteRingtone(c *echo.Context) error {
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	authorID := s.currentUserID(c)
	if authorID == 0 {
		return Render(c, views.OtherError(http.StatusBadRequest, errors.New("You're not logged in.")))
	}

	if err := s.store.DeleteRingtone(c.Request().Context(), displayID, authorID); err != nil {
		return internalError(c, "delete ringtone", err, false)
	}

	soundPath := filepath.Join(s.cfg.RingtonesDir, displayID+".ogg")
	if err := utils.DeleteFile(soundPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to delete ringtone file", "path", soundPath, "display_id", displayID, "error", err)
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (s *Server) detailRingtone(c *echo.Context) error {
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := s.currentUserID(c)

	ringtone, err := s.store.GetRingtone(c.Request().Context(), displayID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.notFound(c)
		}
		return internalError(c, "get ringtone detail", err, true)
	}

	var effects []database.EffectModel
	if userID == 1 {
		effects, err = s.store.GetEffects(c.Request().Context())
		if err != nil {
			return internalError(c, "query effects for ringtone detail", err, true)
		}
	}
	return Render(c, views.Detail(ringtone, effects, userID))
}

func (s *Server) vote(c *echo.Context) error {
	displayID := c.Param("displayID")
	if !validDisplayID(displayID) {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := s.currentUserID(c)
	if userID == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}

	vote, err := strconv.Atoi(c.QueryParam("vote"))
	if err != nil || (vote != 0 && vote != 1 && vote != 2) {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := s.store.Vote(c.Request().Context(), userID, displayID, vote); err != nil {
		slog.Error("failed to update ringtone vote", "method", c.Request().Method, "path", c.Request().URL.Path, "display_id", displayID, "author_id", userID, "error", err)
		ringtone, getErr := s.store.GetRingtone(c.Request().Context(), displayID, userID)
		if getErr != nil {
			return internalError(c, "reload ringtone after failed vote", getErr, false)
		}
		return Render(c, components.Votes(ringtone.DisplayID, ringtone.Votes, ringtone.LoggedInAuthorsVote))
	}

	ringtone, err := s.store.GetRingtone(c.Request().Context(), displayID, userID)
	if err != nil {
		return internalError(c, "reload ringtone after vote", err, false)
	}

	return Render(c, components.Votes(ringtone.DisplayID, ringtone.Votes, ringtone.LoggedInAuthorsVote))
}
