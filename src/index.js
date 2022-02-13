const skinview3d = require("skinview3d");

const canvas = document.getElementById("skindisplay");

let skinViewer = new skinview3d.SkinViewer({
    canvas: canvas,
    width: 600,
    height: 800,
});

// Set the background color
skinViewer.background = 0x5a76f3;


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

document.addEventListener("submit", function (event) {
    event.preventDefault()
    let username = escapeHtml(document.getElementById("username").value)
    let mode = document.querySelector('input[name="mode"]:checked').value;
    fetch(`api/skin/${username}/img/${mode}`)
        .then(response => response.json())
        .then(function (data) {
            let url = data.url;
            if(mode === "full") {
                canvas.style.display = "block";
                skinViewer.loadSkin(url);
            } else {
                canvas.style.display = "none";
            }
            document.getElementById("searchresult").innerHTML = `<img alt='Minecraft skin of ${username}' src='${url}'>`

        })
})
