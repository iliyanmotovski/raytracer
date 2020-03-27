class Particle {
    constructor(light, img) {
        this.light = light;
        this.img = img;
    };

    display() {
        imageMode(CENTER);
        image(this.img, this.light.X, invert(this.light.Y));
    };
}