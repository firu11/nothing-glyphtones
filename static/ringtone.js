import WaveSurfer from "/static/wavesurfer.esm.js"

const listOfRingtones = document.querySelector("#list-of-ringtones")
const all = document.querySelectorAll(".ringtone")
const allWaveSurfers = []

const images = ["/static/icons/play.svg", "/static/icons/pause.svg", "/static/icons/loading.svg"]

function muteAllExcept(index) {
    const imgElements = listOfRingtones.querySelectorAll(".audio button img")
    for (let i = 0; i < allWaveSurfers.length; i++) {
        if (i != index) {
            allWaveSurfers[i].pause()
            if (imgElements[i].getAttribute("src") === images[1]) imgElements[i].src = images[0]
        }
    }
}

for (let i = 0; i < all.length; i++) {
    const id = all[i].getAttribute("data-id")

    const wavesurfer = WaveSurfer.create({
        container: all[i].querySelector(".wave"),
        waveColor: 'white',
        progressColor: 'red',
        url: `/static/sounds/${id}.ogg`,
        barWidth: 5,
        barGap: 5,
        barRadius: 100,
        cursorWidth: 2,
        dragToSeek: true,
        height: 'auto',
        normalize: true,
    })

    wavesurfer.on("ready", () => {
        all[i].querySelector(".audio button img").src = images[0]
    })

    allWaveSurfers.push(wavesurfer)
}

listOfRingtones.addEventListener("click", (e) => {
    if (e.target.tagName == "BUTTON") {
        const i = parseInt(e.target.parentElement.parentElement.getAttribute("data-i"))

        if (e.target.firstChild.getAttribute("src") == images[2]) return

        if (allWaveSurfers[i].isPlaying()) {
            allWaveSurfers[i].pause()
            e.target.firstChild.src = images[0]
        } else {
            muteAllExcept(i)
            allWaveSurfers[i].play()

            e.target.firstChild.src = images[1]
        }
    }
})