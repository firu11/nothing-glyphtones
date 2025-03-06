import {phoneMap} from "./preview-models.js"

let activePhone = null

function showFrame() {
    if (!window.nowPlaying.CSV) return
    const time = window.nowPlaying.player.media.currentTime
    const rowIndex = Math.floor(time * 60)
    
    let glyphRow = window.nowPlaying.CSV[rowIndex]
    if (!glyphRow) { // all zeros
        const len = window.nowPlaying.CSV[0]
        glyphRow = new Array(len)
        for (let i=0; i<len; ++i) a[i] = 0
    }

    phoneMap.get(window.nowPlaying.phoneModel).func(glyphRow)
}

function showPhoneModel(model) {
    phoneMap.forEach(phone => {
        if (model === phone.el.id) phone.el.style.display = "block"
        else phone.el.style.display = "none"
    })
}

function update() {
    if (window.nowPlaying == undefined) window.nowPlaying = {}

    if (window.nowPlaying.isPlaying) {
        //console.log(window.nowPlayingCSV)
        showFrame()
    }
    if (activePhone !== window.nowPlaying.phoneModel) {
        activePhone = window.nowPlaying.phoneModel
        showPhoneModel(window.nowPlaying.phoneModel)
    }
    requestAnimationFrame(update)
}

requestAnimationFrame(update)