const preview = document.getElementById("glyph-preview")
const phoneMap = new Map()
preview.childNodes.forEach(el => {
    if (el === undefined || el.id === undefined) return
    phoneMap.set(el.id, { "el": el })
})

let activePhone = null
const phone = window.innerWidth < 900
const singleRingtonePreview = window.location.pathname.includes("/g/")

function showFrame(rowIndex = 0) {
    if (!window.nowPlaying.CSV) return

    let glyphRow = window.nowPlaying.CSV[rowIndex]
    if (!glyphRow || glyphRow.length == 0) {
        const len = window.nowPlaying.CSV[0].length
        glyphRow = new Array(len).fill("0") // all zeros
    }

    const model = window.nowPlaying.phoneModel
    const glyphs = phoneMap.get(model).el.querySelectorAll("path, rect")
    glyphs.forEach(glyph => {
        const range = glyph.id.split("-")
        if (range.length == 1) { // for simple glyphs
            glyph.setAttribute("opacity", opacity(glyphRow[parseInt(range[0])], 4095))
        } else { // for glyphs split into paths (they need to have a gradient)
            const gradientParts = phoneMap.get(model).el.querySelector("#gradient-" + glyph.id).children
            for (let i = 0; i < gradientParts.length; i++) {
                const colIndex = parseInt(gradientParts[i].id.split("-")[1])
                gradientParts[i].setAttribute("stop-opacity", opacity(glyphRow[colIndex], 4095))
            }
        }
    })
}

function showPhoneModel() {
    if (phone && (activePhone === null || activePhone === undefined)) {
        preview.style.display = "none"
    } else {
        preview.style.display = "block"
    }
    phoneMap.forEach(phone => {
        if (activePhone === phone.el.id) phone.el.style.display = "block"
        else phone.el.style.display = "none"
    })
    if (["(1)", "(1)_15", "(2)"].includes(activePhone)) {
        preview.classList.replace("square", "rectangle")
        preview.classList.replace("circle", "rectangle")
    } else if (["(2a)", "(3a)"].includes(activePhone)) {
        preview.classList.replace("rectangle", "square")
        preview.classList.replace("circle", "square")
    } else if (["(3)"].includes(activePhone)) {
        preview.classList.replace("rectangle", "circle")
        preview.classList.replace("square", "circle")
    }
}

function update() {
    if (window.nowPlaying == undefined) window.nowPlaying = {}

    if (window.nowPlaying.isPlaying) {
        // console.log(window.nowPlaying.CSV)
        const time = window.nowPlaying.player.media.currentTime
        const rowIndex = Math.ceil(time * 60)
        showFrame(rowIndex)
    }
    if (activePhone !== window.nowPlaying.phoneModel) {
        activePhone = window.nowPlaying.phoneModel
        showPhoneModel()
    }
    requestAnimationFrame(update)
}

const minOpacity = 0.1
function opacity(value, maxValue = 4095) {
    return minOpacity + (1 - minOpacity) * (value / maxValue)
}

requestAnimationFrame(update)
