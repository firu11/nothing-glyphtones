document.querySelector('.share-btn').addEventListener('click', async () => {
    const url = window.location.href
    try {
        await navigator.share({
            title: document.title,
            text: 'Check out my Glyphtone!',
            url,
        })
    } catch (err) {
        alert("Whoops, share canceled or failed. Your browser might not support this feature.\n\nYou can just copy the link in your adress bar and send it yourself.")
        console.log("Share canceled or failed: " + err)
    }
})
