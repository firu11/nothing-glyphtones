const preview = document.getElementById("glyph-preview")
let initialX = 0, initialY = 0
let currentTranslateX = 0, currentTranslateY = 0
let isDragging = false

// Get the element's default position
const defaultRect = preview.getBoundingClientRect()
const defaultX = defaultRect.left
const defaultY = defaultRect.top

// Start dragging (mouse or touch)
function startDragging(e) {
    e.preventDefault()
    if (window.getComputedStyle(preview).position !== "fixed") return

    isDragging = true

    // Handle both mouse and touch
    const event = e.type === "touchstart" ? e.touches[0] : e
    initialX = event.clientX
    initialY = event.clientY

    const transform = preview.style.transform
    if (transform) {
        const match = transform.match(/translate\(([^,]+)px, ([^)]+)px\)/)
        if (match) {
            currentTranslateX = parseFloat(match[1])
            currentTranslateY = parseFloat(match[2])
        }
    }

    // Add move and stop listeners only when dragging starts
    document.addEventListener("mousemove", moveDragging)
    document.addEventListener("mouseup", stopDragging)
    document.addEventListener("touchmove", moveDragging, { passive: false })
    document.addEventListener("touchend", stopDragging)
}

// Move element (mouse or touch)
function moveDragging(e) {
    if (!isDragging) return
    e.preventDefault()

    // Handle both mouse and touch
    const event = e.type === "touchmove" ? e.touches[0] : e
    const deltaX = event.clientX - initialX
    const deltaY = event.clientY - initialY

    // Calculate new position
    let newTranslateX = currentTranslateX + deltaX
    let newTranslateY = currentTranslateY + deltaY

    // Get screen and element dimensions
    const screenWidth = window.innerWidth
    const screenHeight = window.innerHeight
    const elementWidth = preview.offsetWidth
    const elementHeight = preview.offsetHeight

    // Calculate boundaries relative to the default position
    const minX = -defaultX
    const maxX = screenWidth - (defaultX + elementWidth)
    const minY = -defaultY
    const maxY = screenHeight - (defaultY + elementHeight)

    // Clamp position within boundaries
    newTranslateX = Math.max(minX, Math.min(newTranslateX, maxX))
    newTranslateY = Math.max(minY, Math.min(newTranslateY, maxY))

    // Apply new position
    preview.style.transform = `translate(${newTranslateX}px, ${newTranslateY}px)`
}

// Stop dragging (mouse or touch)
function stopDragging(e) {
    e.preventDefault()
    isDragging = false

    // Remove move and stop listeners when dragging stops
    document.removeEventListener("mousemove", moveDragging)
    document.removeEventListener("mouseup", stopDragging)
    document.removeEventListener("touchmove", moveDragging)
    document.removeEventListener("touchend", stopDragging)
}

// Mouse event listeners
preview.addEventListener("mousedown", startDragging)

// Touch event listeners
preview.addEventListener("touchstart", startDragging)