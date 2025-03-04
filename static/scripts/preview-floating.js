const preview = document.getElementById("glyph-preview")
let initialMouseX = 0, initialMouseY = 0
let currentTranslateX = 0, currentTranslateY = 0

// Get the element's default position
const defaultRect = preview.getBoundingClientRect()
const defaultX = defaultRect.left
const defaultY = defaultRect.top

function startMoving(e) {
    e.preventDefault()
    if (window.getComputedStyle(preview).position !== "fixed") return
    initialMouseX = e.clientX
    initialMouseY = e.clientY
    const transform = preview.style.transform
    if (transform) {
        const match = transform.match(/translate\(([^,]+)px, ([^)]+)px\)/)
        if (match) {
            currentTranslateX = parseFloat(match[1])
            currentTranslateY = parseFloat(match[2])
        }
    }
    document.addEventListener("mousemove", move)
    document.addEventListener("mouseup", stopMoving)
}

function stopMoving(e) {
    e.preventDefault()
    document.removeEventListener("mousemove", move)
    document.removeEventListener("mouseup", stopMoving)
}

function move(e) {
    e.preventDefault()
    const deltaX = e.clientX - initialMouseX
    const deltaY = e.clientY - initialMouseY

    // Calculate new position
    let newTranslateX = currentTranslateX + deltaX
    let newTranslateY = currentTranslateY + deltaY

    // Get the screen dimensions
    const screenWidth = window.innerWidth
    const screenHeight = window.innerHeight

    // Get the element dimensions
    const elementWidth = preview.offsetWidth
    const elementHeight = preview.offsetHeight

    // Calculate the boundaries relative to the default position
    const minX = -defaultX // Prevent moving left beyond the default position
    const maxX = screenWidth - (defaultX + elementWidth) // Prevent moving right beyond the screen
    const minY = -defaultY // Prevent moving up beyond the default position
    const maxY = screenHeight - (defaultY + elementHeight) // Prevent moving down beyond the screen

    // Clamp the new position within the boundaries
    newTranslateX = Math.max(minX, Math.min(newTranslateX, maxX))
    newTranslateY = Math.max(minY, Math.min(newTranslateY, maxY))

    // Apply the new position
    preview.style.transform = `translate(${newTranslateX}px, ${newTranslateY}px)`
}

preview.addEventListener("mousedown", startMoving)