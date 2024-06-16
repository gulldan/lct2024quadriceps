"use client"
const config = {
  heat: 0,
  colors: { start: `rgb(255, 255, 0)`, end: `rgb(255, 150, 0)` },
};

function radomMinMax(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}
function randomXY(cx, cy) {
  return {
    x: radomMinMax(this.x - cx, this.x + cx),
    y: radomMinMax(this.y - cy * Math.random(), this.y - cy * Math.random()),
  };
}

class Particle {
  constructor(fireplace) {
    this.randomXY = randomXY.bind(this);
    this.fireplace = fireplace;
    this.x = radomMinMax(
      this.fireplace.size.cx - 50,
      this.fireplace.size.cx + 50
    );
    this.y = this.fireplace.size.cy;

    this.size = 8;
    this.color = config.colors.start;
    this.bezierPoints = [
      { x: this.x, y: this.y },
      this.randomXY(100, 50),
      this.randomXY(80, 150),
      this.randomXY(10, 300),
    ];
    this.speed = 0.01;
    this.t = 0;
  }
  updateOnMouse() {
    let dx = this.fireplace.mouse.x - this.x;
    let dy = this.fireplace.mouse.y - this.y;
    let distance = dx ** 2 + dy ** 2;
    let force = -this.fireplace.mouse.radius / distance;
    let angle = 0;

    if (distance < this.fireplace.mouse.radius) {
      angle = Math.atan2(dy, dx);
      this.bezierPoints.forEach((point) => {
        point.x += force * Math.cos(angle);
        point.y += force * Math.sin(angle);
      });
    }
  }
  updateColors() {
    let [, r, g, b] = this.fireplace.rgb.start;
    let dx = this.fireplace.size.cx - this.x;
    let dy = this.fireplace.size.cy - this.y;
    let distance = dx ** 2 + dy ** 2;
    r = Math.ceil(
      Math.max(this.fireplace.rgb.end[1], r - distance * this.speed)
    );
    g = Math.ceil(
      Math.max(this.fireplace.rgb.end[2], g - distance * this.speed)
    );
    b = Math.ceil(
      Math.max(this.fireplace.rgb.end[3], b - distance * this.speed)
    );
    this.color = `rgb(${[r, g, b].join(",")})`;
  }
  updateParticles([p0, p1, p2, p3]) {
    // * Calculate coefficients based on a particle current position
    let cx = 3 * (p1.x - p0.x);
    let bx = 3 * (p2.x - p1.x) - cx;
    let ax = p3.x - p0.x - cx - bx;

    let cy = 3 * (p1.y - p0.y);
    let by = 3 * (p2.y - p1.y) - cy;
    let ay = p3.y - p0.y - cy - by;

    this.t += this.speed;
    // * Calculate new X & Y positions
    let xt =
      ax * (this.t * this.t * this.t) +
      bx * (this.t * this.t) +
      cx * this.t +
      p0.x;
    let yt =
      ay * (this.t * this.t * this.t) +
      by * (this.t * this.t) +
      cy * this.t +
      p0.y;

    if (this.t > 1) this.t = 0;

    this.size -= 0.05;
    if (this.size < 0.5) this.size = 0.5;

    this.x = xt;
    this.y = yt;
  }

  update() {
    this.updateParticles(this.bezierPoints);
    this.updateOnMouse();
    this.updateColors();
  }
  draw(context) {
    context.fillStyle = this.color;
    context.fillRect(this.x, this.y, this.size, this.size);
  }
}

export class Fireplace {
  constructor() {
    this.cnv = null;
    this.ctx = null;
    this.size = { w: 0, h: 0, cx: 0, cy: 0 };
    this.particles = [];
    this.rgb = {
      start: /rgb\((\d{1,3}), (\d{1,3}), (\d{1,3})\)/.exec(config.colors.start),
      end: /rgb\((\d{1,3}), (\d{1,3}), (\d{1,3})\)/.exec(config.colors.end),
    };
    this.particlesSpawnRate = 10;
    this.mouse = {
      radius: 3000,
      x: undefined,
      y: undefined,
    };
    this.config = config
    window.addEventListener("mousemove", (event) => {
      this.mouse.x = event.x;
      this.mouse.y = event.y;
    });
    window.addEventListener("touchmove", (event) => {
      this.mouse.x = event.touches[0].clientX;
      this.mouse.y = event.touches[0].clientY;
    });
  }
  init() {
    this.inited = true;
    this.createCanvas();
    this.updateAnimation();
  }
  createCanvas() {
    this.cnv = document.getElementById(`cvs`);
    if (!this.cnv) {
        console.error("cnv is none")
        return
    }
    this.ctx = this.cnv.getContext(`2d`);
    this.setCanvasSize();
    window.addEventListener(`resize`, () => this.setCanvasSize());
  }
  setCanvasSize() {
    this.size.w = this.cnv.width = window.innerWidth;
    this.size.h = this.cnv.height = window.innerHeight;
    this.size.cx = this.size.w / 2;
    this.size.cy = this.size.h / 2 + 200;
  }
  generateParticles() {
    for (let i = 0; i < this.particlesSpawnRate; i++) {
      this.particles.push(new Particle(this));
    }
    let particlesShift =
      this.particles.length > this.config.heat
        ? this.particlesSpawnRate
        : this.particlesSpawnRate / 2;
    for (let i = 0; i < particlesShift; i++) {
      this.particles.shift();
    }
  }
  drawParticles() {
    this.particles.forEach((particle) => particle.update());
    this.particles.forEach((particle) => particle.draw(this.ctx));
    this.generateParticles();
  }

  updateCavas() {
    this.ctx.fillStyle = `rgb(22, 22, 25)`;
    this.ctx.fillRect(0, 0, this.size.w, this.size.h);
    this.ctx.shadowColor = this.config.colors.end;
    this.ctx.shadowBlur = 25;
  }
  updateAnimation() {
    this.updateCavas();
    this.drawParticles();
    requestAnimationFrame(() => this.updateAnimation());
  }
}
