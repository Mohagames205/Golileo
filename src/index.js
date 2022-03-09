const skinview3d = require("skinview3d");

const canvas = document.getElementById("skindisplay");

let skinViewer = new skinview3d.SkinViewer({
    canvas: canvas,
    width: 600,
    height: 800,
    preserveDrawingBuffer: true
});

// Change camera FOV
skinViewer.fov = 70;

// Zoom out
skinViewer.zoom = 0.5;

// Control objects with your mouse!
let control = skinview3d.createOrbitControls(skinViewer);
control.enableRotate = true;
control.enableZoom = false;
control.enablePan = false;


skinViewer.animations.add(skinview3d.IdleAnimation);

function escapeHtml(text) {
    var map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };

    return text.replace(/[&<>"']/g, function(m) { return map[m]; });
}

const themeToggle = document.getElementById('theme-toggle');
const body = document.body;

themeToggle.onclick = () => {
    if(body.classList.contains('light')){
        body.classList.replace('light', 'dark');
        themeToggle.innerText = "Light mode ☀️"
    } else if(body.classList.contains('dark')){
        body.classList.replace('dark', 'light');
        themeToggle.innerText = "Dark mode 🌜"
    }
}

document.addEventListener("submit", function (event) {
    event.preventDefault()
    let username = escapeHtml(document.getElementById("username").value)
    let mode = document.querySelector('input[name="mode"]:checked').value;
    let url = `api/${username}/img/${mode}`

    if(mode === "full") {
        canvas.style.display = "block";
        skinViewer.loadSkin(url);
    } else {
        canvas.style.display = "none";
    }
    document.getElementById("searchresult").innerHTML = `<img alt='Minecraft skin of ${username}' src='${url}'>`
})
