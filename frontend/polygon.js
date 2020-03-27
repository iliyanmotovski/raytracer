class Polygons {
    constructor(polygons, color, stroke) {
        this.polygons = polygons;
        this.color = color;
        this.stroke = stroke;
    };

    display() {
        this.polygons.map(polygon => {
            fill(this.color[0], this.color[1], this.color[2]);
            if (!this.stroke) {
                noStroke();
            }

            beginShape();
            polygon.forEach(vertice => vertex(vertice.X, invert(vertice.Y)));
            endShape(CLOSE);
        });
    };
}