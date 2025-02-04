import WaveSurfer from "/static/scripts/wavesurfer.esm.js"

const images = ["/static/icons/play.svg", "/static/icons/pause.svg", "/static/icons/loading.svg"]
const imagesRed = ["/static/icons/play-red.svg", "/static/icons/pause-red.svg", "/static/icons/loading.svg"]

let form, fileInputElement, fileInputImageElement, fileInputSpanElement, fileInputAudioPreviewElement, fileInputAudioPlayElement, wave
let wavesurfer

function handleFile(e) {
    let file = e.srcElement.files[0]
    if (file == undefined) {
        fileInputImageElement.style.display = "block"
        fileInputSpanElement.innerText = "Click to choose a file."
        fileInputAudioPreviewElement.style.display = "none"
    } else {
        fileInputSpanElement.innerHTML = file.name + ' <img src="/static/icons/edit.svg" width="16"/>'
        fileInputImageElement.style.display = "none"
        fileInputAudioPreviewElement.style.display = "flex"

        audio(file)
    }
}

function formChange(e) {
    const errorMsg = form.querySelector("h3.red-heading")
    if (errorMsg !== null) errorMsg.remove()
}

function audio(file) {
    if (wavesurfer !== undefined) wavesurfer.destroy()
    wavesurfer = WaveSurfer.create({
        container: wave,
        waveColor: 'white',
        progressColor: 'red',
        url: URL.createObjectURL(file),
        barWidth: 4,
        barGap: 4,
        barRadius: 100,
        cursorWidth: 2,
        dragToSeek: true,
        height: 'auto',
        normalize: true,
    })

    wavesurfer.on("ready", () => {
        fileInputAudioPreviewElement.querySelector("button img.white").src = images[0]
        fileInputAudioPreviewElement.querySelector("button img.red").src = imagesRed[0]
    })

    wavesurfer.on("finish", () => {
        fileInputAudioPreviewElement.querySelector("button img.white").src = images[0]
        fileInputAudioPreviewElement.querySelector("button img.red").src = imagesRed[0]
    })
}

function click(e) {
    if (e.target.firstChild.getAttribute("src") == images[2]) return

    if (wavesurfer.isPlaying()) {
        wavesurfer.pause()
        e.target.querySelector(".white").src = images[0]
        e.target.querySelector(".red").src = imagesRed[0]
    } else {
        wavesurfer.play()
        e.target.querySelector(".white").src = images[1]
        e.target.querySelector(".red").src = imagesRed[1]
    }
}

function reloadSelectors() {
    form = document.querySelector("form#upload")
    if (form == null) {
        return
    }

    fileInputElement = form.querySelector("input[type='file']")
    fileInputImageElement = form.querySelector("#image")
    fileInputSpanElement = form.querySelector("span")
    fileInputAudioPreviewElement = form.querySelector("#audio")
    fileInputAudioPlayElement = fileInputAudioPreviewElement.querySelector("button")
    wave = fileInputAudioPreviewElement.querySelector(".wave")

    form.addEventListener("change", formChange)
    fileInputElement.addEventListener("change", handleFile)
    fileInputAudioPlayElement.addEventListener("click", click)
}

reloadSelectors()
document.addEventListener("htmx:afterSettle", reloadSelectors)