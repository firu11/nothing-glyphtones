const preview = document.getElementById("glyph-preview")

const funcMap = new Map([
    ["(1)", phone1],
    ["(2)", phone2],
    ["(2a)", phone2a],
])
export const phoneMap = new Map()
preview.childNodes.forEach(el => {
    phoneMap.set(el.id, { "el": el, "func": funcMap.get(el.id) })
})

export function phone1(csvRow) {
    const glyphs = phoneMap.get("(1)").el.querySelectorAll("path")
    glyphs.forEach(glyph => {
        const range = glyph.id.split("-")
        if (range.length == 1) {
            glyph.setAttribute("opacity", opacity(csvRow[parseInt(range[0])], 4095))
        } else {
            const numberOfParts = range[1] - range[0] + 1
            let sum = 0
            for (let i = parseInt(range[0]); i <= parseInt(range[1]); i++) {
                sum += parseInt(csvRow[i])
            }
            glyph.setAttribute("opacity", opacity(sum, numberOfParts * 4095))
        }
    })
}

export function phone2and2a(csvRow, id) {
    const glyphs = phoneMap.get(id).el.querySelectorAll("path")
    glyphs.forEach(glyph => {
        const range = glyph.id.split("-")
        if (range.length == 1) { // for simple glyphs
            glyph.setAttribute("opacity", opacity(csvRow[parseInt(range[0])], 4095))
        } else { // for glyphs split into paths (they need to have a gradient)
            const gradientParts = phoneMap.get(id).el.querySelector("#gradient-" + glyph.id).children
            for (let i = 0; i < gradientParts.length; i++) {
                const colIndex = parseInt(gradientParts[i].id.split("-")[1])
                gradientParts[i].setAttribute("stop-opacity", opacity(csvRow[colIndex], 4095))
            }
        }
    })
}

const minOpacity = 0.1
function opacity(value, maxValue = 4095) {
    return minOpacity + (1 - minOpacity) * (value / maxValue)
}