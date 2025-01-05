import WaveSurfer from "/static/scripts/wavesurfer.esm.js"

const images = ["/static/icons/play.svg", "/static/icons/pause.svg", "/static/icons/loading.svg"]
const imagesRed = ["/static/icons/play-red.svg", "/static/icons/pause-red.svg", "/static/icons/loading.svg"]
let allWaveSurfers = []
let listOfRingtones = []
let all = []

function muteAllExcept(index) {
    const imgElements = listOfRingtones.querySelectorAll(".audio button img.white")
    const imgRedElements = listOfRingtones.querySelectorAll(".audio button img.red")
    for (let i = 0; i < allWaveSurfers.length; i++) {
        if (i != index) {
            allWaveSurfers[i].pause()
            if (imgElements[i].getAttribute("src") === images[1]) imgElements[i].src = images[0]
            if (imgRedElements[i].getAttribute("src") === imagesRed[1]) imgRedElements[i].src = imagesRed[0]
        }
    }
}

function click(e) {
    if (e.target.tagName == "BUTTON" && e.target.firstChild.tagName == "IMG") {
        const i = parseInt(e.target.parentElement.parentElement.getAttribute("data-i"))

        if (e.target.firstChild.getAttribute("src") == images[2]) return

        if (allWaveSurfers[i].isPlaying()) {
            allWaveSurfers[i].pause()
            e.target.querySelector(".white").src = images[0]
            e.target.querySelector(".red").src = imagesRed[0]
        } else {
            muteAllExcept(i)
            allWaveSurfers[i].play()

            e.target.querySelector(".white").src = images[1]
            e.target.querySelector(".red").src = imagesRed[1]
        }
    }
}

function main(e) {
    if (e !== undefined && e.detail.elt.id !== "list-of-ringtones") return // only if the target is list of ringtones

    document.removeEventListener("click", click)

    listOfRingtones = document.querySelector("#list-of-ringtones")
    all = document.querySelectorAll(".ringtone")
    allWaveSurfers = []

    for (let i = 0; i < all.length; i++) {
        const id = all[i].getAttribute("data-id")

        const wavesurfer = WaveSurfer.create({
            container: all[i].querySelector(".wave"),
            waveColor: 'white',
            progressColor: 'red',
            url: `/sounds/${id}.ogg`,
            barWidth: 4,
            barGap: 4,
            barRadius: 100,
            cursorWidth: 2,
            dragToSeek: true,
            height: 'auto',
            normalize: true,
        })

        wavesurfer.on("ready", () => {
            all[i].querySelector(".audio button img.white").src = images[0]
            all[i].querySelector(".audio button img.red").src = imagesRed[0]
        })

        wavesurfer.on("finish", () => {
            all[i].querySelector(".audio button img.white").src = images[0]
            all[i].querySelector(".audio button img.red").src = imagesRed[0]
        })

        allWaveSurfers.push(wavesurfer)
    }

    listOfRingtones.addEventListener("click", click)
}

main()
document.addEventListener("htmx:afterSwap", main)