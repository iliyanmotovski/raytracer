let img;
let scene;
let columns;
let particle;
let triangles;
let refreshIntervalId;

let dragged = false;
let getSceneUrl = 'http://localhost:8008/api/v1/scene'
let postConfigUrl = 'http://localhost:8008/api/v1/scene/config'

function preload() {
    httpGet(getSceneUrl, 'json', false, resp, err);
}

function setupScene() {
    img = loadImage('sun.png');
    createCanvas(scene.Width, scene.Height);
}

function draw() {
    if (!scene) { return }
    background(117, 114, 107);

    columns = new Polygons(scene.Polygons, [181, 121, 24], true);
    columns.display();

    triangles = new Polygons(scene.Triangles, [217, 206, 189], false);
    triangles.display();

    fill(0,0,0);
    textSize(19);
    text('Lit area is: ' + scene.LitArea + '%', 10, 30);

    particle = new Particle(scene.Light, img);

    if (dragged) {
        scene.Light.X = mouseX;
        scene.Light.Y = invert(mouseY);
    }

    particle.display();
}

function mousePressed() {
    let x = scene.Light.X;
    let y = invert(scene.Light.Y);
    let width = img.width;
    let height = img.height;

    if (mouseX > x && mouseX < x + width && mouseY > y && mouseY < y + height) {
        dragged = true;
        updateConfig(100)
    }
}

function mouseReleased() {
    dragged = false;
    clearInterval(refreshIntervalId)
}

function updateConfig(interval) {
    refreshIntervalId = setInterval(() => {
        postData = {light: scene.Light, polygons: scene.Polygons, scene: {X: scene.Width, Y: scene.Height}};
        httpPost(postConfigUrl, 'json', postData, () => {
            httpGet(getSceneUrl, 'json', false, resp, err);
        }, err);
    },interval);
}

function resp(response) {
    scene = response;
    setupScene();
}

function err(err) {
    textSize(32);
    text(err, 10, 30);
    fill(0, 102, 153);
}