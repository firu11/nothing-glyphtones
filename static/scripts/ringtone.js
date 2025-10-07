import WaveSurfer from "/static/scripts/wavesurfer.esm.js"
import Pako from "/static/scripts/pako.esm.mjs"

const images = ["/static/icons/play.svg", "/static/icons/pause.svg", "/static/icons/loading.svg"]
const imagesRed = ["/static/icons/play-red.svg", "/static/icons/pause-red.svg", "/static/icons/loading.svg"]
let allWaveSurfers = []
let listOfRingtones = []
let all = []
const singleRingtonePreview = window.location.pathname.includes("/g/")

function muteAllExcept(index) {
    if (allWaveSurfers.length === 1) return
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
        const ringtoneDiv = e.target.parentElement.parentElement
        const i = parseInt(ringtoneDiv.getAttribute("data-i"))

        if (e.target.firstChild.getAttribute("src") == images[2]) return

        if (allWaveSurfers[i].isPlaying()) {
            allWaveSurfers[i].pause()
            if (!singleRingtonePreview)
                window.nowPlaying.phoneModel = null
            window.nowPlaying.player = null
            window.nowPlaying.isPlaying = false
            e.target.querySelector(".white").src = images[0]
            e.target.querySelector(".red").src = imagesRed[0]
        } else {
            muteAllExcept(i)
            allWaveSurfers[i].play()
            window.nowPlaying.player = allWaveSurfers[i]
            setPhoneModel(ringtoneDiv)

            const glyphs = ringtoneDiv.getAttribute("data-glyphs")
            let resultCSV = ""
            try {
                const compressedData = atob(glyphs)

                const bytes = new Uint8Array(compressedData.length)
                for (let i = 0; i < compressedData.length; i++) {
                    bytes[i] = compressedData.charCodeAt(i)
                }
                resultCSV = Pako.inflate(bytes, { to: 'string' })
            } catch (err) {
                console.log(err)
            }

            if (resultCSV !== undefined) {
                let rows = resultCSV.split(/\r\n|\n/)
                let csv = []
                rows.forEach(row => {
                    csv.push(row.split(",").slice(0, -1))
                })
                window.nowPlaying.CSV = csv
            }

            window.nowPlaying.isPlaying = true
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

    window.nowPlaying = {}

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
            window.nowPlaying.CSV = ""
            window.nowPlaying.isPlaying = false
            if (!singleRingtonePreview)
                window.nowPlaying.phoneModel = null
            window.nowPlaying.player = null
        })

        allWaveSurfers.push(wavesurfer)
    }

    if (listOfRingtones === null) {
        all[0].addEventListener("click", click)
        setPhoneModel(all[0])
    } else {
        listOfRingtones.addEventListener("click", click)
    }
}

function setPhoneModel(ringtoneDiv) {
    const phones = ringtoneDiv.getAttribute("data-phone").split(",")
    if (phones.length == 1 && phones[0] == "(1)") { // 15 zone (1)
        window.nowPlaying.phoneModel = "(1)_15"
    } else window.nowPlaying.phoneModel = phones[0]
}

main()
document.addEventListener("htmx:afterSwap", main)

document.body.addEventListener("htmx:responseError", function (event) {
    if (event.detail.xhr.status === 401) {
        const messageBox = document.getElementById("unauthorized-message")
        if (messageBox === null) return
        messageBox.innerText = "Unauthorized! Please log in."
        messageBox.style.display = "block"
        setInterval(() => {
            messageBox.style.display = "none"
        }, 4000)
    }
})
