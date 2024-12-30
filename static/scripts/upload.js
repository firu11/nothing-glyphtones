let form, fileInputElement, fileInputImageElement, fileInputSpanElement

function handleFile(e) {
    let file = e.srcElement.files[0]
    if (file == undefined) {
        fileInputImageElement.style.display = "block"
        fileInputSpanElement.style.display = "none"
    } else {
        fileInputSpanElement.innerText = file.name

        fileInputImageElement.style.display = "none"
        fileInputSpanElement.style.display = "block"
    }
}

function formChange(e) {
    const errorMsg = form.querySelector("h3.red-heading")
    if (errorMsg !== null) errorMsg.remove()
}

function reloadSelectors() {
    form = document.querySelector("form#upload")
    if (form == null) {
        return
    }
    
    fileInputElement = form.querySelector("input[type='file']")
    fileInputImageElement = form.querySelector("img")
    fileInputSpanElement = form.querySelector("span")

    form.addEventListener("change", formChange)
    fileInputElement.addEventListener("change", handleFile)

}

reloadSelectors()
document.addEventListener("htmx:afterSettle", reloadSelectors)