const fileInputLabelElement = document.querySelector("label:has(input[type='file'])")
const fileInputElement = fileInputLabelElement.querySelector("input[type='file']")
const fileInputImageElement = fileInputLabelElement.querySelector("img")
const fileInputSpanElement = fileInputLabelElement.querySelector("span")
fileInputElement.addEventListener("change", handleFile)

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
    console.log(file)
}