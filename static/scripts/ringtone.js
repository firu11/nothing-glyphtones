import WaveSurfer from '/static/scripts/wavesurfer.esm.js';
import Pako from '/static/scripts/pako.esm.mjs';

const images = ['/static/icons/play.svg', '/static/icons/pause.svg', '/static/icons/loading.svg'];
const imagesRed = ['/static/icons/play-red.svg', '/static/icons/pause-red.svg', '/static/icons/loading.svg'];
let allWaveSurfers = [];
let listOfRingtones = [];
let all = [];
let playRequest = 0;
const glyphCache = new Map();
const singleRingtonePreview = window.location.pathname.includes('/g/');

function muteAllExcept(index) {
  if (allWaveSurfers.length === 1) return;
  const imgElements = listOfRingtones.querySelectorAll('.audio button img.white');
  const imgRedElements = listOfRingtones.querySelectorAll('.audio button img.red');
  for (let i = 0; i < allWaveSurfers.length; i++) {
    if (i != index) {
      allWaveSurfers[i].pause();
      if (imgElements[i].getAttribute('src') === images[1]) imgElements[i].src = images[0];
      if (imgRedElements[i].getAttribute('src') === imagesRed[1]) imgRedElements[i].src = imagesRed[0];
    }
  }
}

function setButtonState(button, state) {
  button.querySelector('.white').src = images[state];
  button.querySelector('.red').src = imagesRed[state];
}

function inflateGlyphs(buffer) {
  const bytes = new Uint8Array(buffer);
  try {
    return Pako.inflate(bytes, { to: 'string' });
  } catch {
    const base64 = new TextDecoder().decode(bytes).replace(/\s/g, '');
    const compressed = Uint8Array.from(atob(base64), (character) => character.charCodeAt(0));
    return Pako.inflate(compressed, { to: 'string' });
  }
}

function loadGlyphCSV(id) {
  if (!glyphCache.has(id)) {
    const request = fetch(`/ringtones/${id}/glyphs`)
      .then((response) => {
        if (!response.ok) throw new Error(`Failed to load glyphs: ${response.status}`);
        return response.arrayBuffer();
      })
      .then(inflateGlyphs)
      .then((value) => value.split(/\r?\n/).map((row) => row.split(',').slice(0, -1)))
      .catch((error) => {
        glyphCache.delete(id);
        throw error;
      });
    glyphCache.set(id, request);
  }
  return glyphCache.get(id);
}

async function click(event) {
  const button = event.target.closest('.audio button');
  if (!button || button.querySelector('.white').getAttribute('src') === images[2]) return;

  const ringtoneDiv = button.closest('.ringtone');
  const index = Number(ringtoneDiv.dataset.i);
  const player = allWaveSurfers[index];

  if (player.isPlaying()) {
    playRequest++;
    player.pause();
    window.nowPlaying = singleRingtonePreview ? { phoneModel: window.nowPlaying.phoneModel } : {};
    setButtonState(button, 0);
    return;
  }

  const request = ++playRequest;
  const cancelled = () => request !== playRequest || !ringtoneDiv.isConnected;
  window.nowPlaying = {};
  setButtonState(button, 2);

  const glyphRequest = loadGlyphCSV(ringtoneDiv.dataset.id).catch((error) => {
    console.error(error);
    return null;
  });

  muteAllExcept(index);
  try {
    await player.play();
  } catch (error) {
    console.error(error);
    setButtonState(button, 0);
    return;
  }

  if (cancelled()) {
    player.pause();
    if (ringtoneDiv.isConnected) setButtonState(button, 0);
    return;
  }

  window.nowPlaying = { CSV: null, player: player, isPlaying: true };
  setPhoneModel(ringtoneDiv);
  setButtonState(button, 1);

  const csv = await glyphRequest;
  if (csv !== null && !cancelled()) window.nowPlaying.CSV = csv;
}

function main(e) {
  if (e !== undefined && e.detail.elt.id !== 'list-of-ringtones') return; // only if the target is list of ringtones

  listOfRingtones = document.querySelector('#list-of-ringtones');
  all = document.querySelectorAll('.ringtone');
  allWaveSurfers.forEach((wavesurfer) => wavesurfer.destroy());
  allWaveSurfers = [];
  playRequest++;

  window.nowPlaying = {};

  for (let i = 0; i < all.length; i++) {
    const id = all[i].getAttribute('data-id');

    const wavesurfer = WaveSurfer.create({
      container: all[i].querySelector('.wave'),
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
    });

    wavesurfer.on('ready', () => {
      all[i].querySelector('.audio button img.white').src = images[0];
      all[i].querySelector('.audio button img.red').src = imagesRed[0];
    });

    wavesurfer.on('finish', () => {
      all[i].querySelector('.audio button img.white').src = images[0];
      all[i].querySelector('.audio button img.red').src = imagesRed[0];
      window.nowPlaying.CSV = '';
      window.nowPlaying.isPlaying = false;
      if (!singleRingtonePreview) window.nowPlaying.phoneModel = null;
      window.nowPlaying.player = null;
    });

    allWaveSurfers.push(wavesurfer);
  }

  if (listOfRingtones === null) {
    all[0].addEventListener('click', click);
    setPhoneModel(all[0]);
  } else {
    listOfRingtones.addEventListener('click', click);
  }
}

function setPhoneModel(ringtoneDiv) {
  const phones = ringtoneDiv.getAttribute('data-phone').split(',');
  if (phones.length == 1 && phones[0] == '(1)') {
    // 15 zone (1)
    window.nowPlaying.phoneModel = '(1)_15';
  } else window.nowPlaying.phoneModel = phones[0];
}

main();
document.addEventListener('htmx:afterSwap', main);

document.body.addEventListener('htmx:responseError', function (event) {
  if (event.detail.xhr.status === 401) {
    const messageBox = document.getElementById('unauthorized-message');
    if (messageBox === null) return;
    messageBox.innerText = 'Unauthorized! Please log in.';
    messageBox.style.display = 'block';
    setInterval(() => {
      messageBox.style.display = 'none';
    }, 4000);
  }
});
