// chatgpt made this shit

const preview = document.getElementById("glyph-preview");

let isDragging = false;
let initialX = 0, initialY = 0;
let startTranslateX = 0, startTranslateY = 0;

function getTranslate() {
    const style = window.getComputedStyle(preview);
    const transform = style.transform;
    if (!transform || transform === "none") return { x: 0, y: 0 };

    const match = transform.match(/matrix\((.+)\)/);
    if (match) {
        const parts = match[1].split(", ");
        return {
            x: parseFloat(parts[4]),
            y: parseFloat(parts[5]),
        };
    }

    return { x: 0, y: 0 };
}

function startDragging(e) {
    if (e.type === "mousedown" && e.button !== 0) return; // only left-click
    e.preventDefault();

    const event = e.type.startsWith("touch") ? e.touches[0] : e;
    isDragging = true;

    initialX = event.clientX;
    initialY = event.clientY;

    const { x, y } = getTranslate();
    startTranslateX = x;
    startTranslateY = y;

    document.addEventListener("mousemove", moveDragging);
    document.addEventListener("mouseup", stopDragging);
    document.addEventListener("touchmove", moveDragging, { passive: false });
    document.addEventListener("touchend", stopDragging);
}

function moveDragging(e) {
    if (!isDragging) return;
    e.preventDefault();

    const event = e.type.startsWith("touch") ? e.touches[0] : e;
    const deltaX = event.clientX - initialX;
    const deltaY = event.clientY - initialY;

    let translateX = startTranslateX + deltaX;
    let translateY = startTranslateY + deltaY;

    // Get element's live size
    const width = preview.offsetWidth;
    const height = preview.offsetHeight;

    // Get fixed left/top position from computed style
    const style = window.getComputedStyle(preview);
    const fixedLeft = parseFloat(style.left);
    const fixedTop = parseFloat(style.top);

    // Clamp to window
    const minX = -fixedLeft;
    const maxX = window.innerWidth - fixedLeft - width;
    const minY = -fixedTop;
    const maxY = window.innerHeight - fixedTop - height;

    translateX = Math.max(minX, Math.min(translateX, maxX));
    translateY = Math.max(minY, Math.min(translateY, maxY));

    preview.style.transform = `translate(${translateX}px, ${translateY}px)`;
}

function stopDragging(e) {
    isDragging = false;

    document.removeEventListener("mousemove", moveDragging);
    document.removeEventListener("mouseup", stopDragging);
    document.removeEventListener("touchmove", moveDragging);
    document.removeEventListener("touchend", stopDragging);
}

preview.addEventListener("mousedown", startDragging);
preview.addEventListener("touchstart", startDragging, { passive: false });
