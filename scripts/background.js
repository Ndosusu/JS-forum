const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');

// Valeurs fixes pour les paramètres des vagues
const WAVE_SPEED = 10 / 20000; 
const WAVE_AMPLITUDE =70;    
const WAVE_DENSITY = 8;       
const WAVE_COLOR = "#da9052"; 

let waves = [];
let mouseX = 0;
let mouseY = 0;

// Suivre la position de la souris
document.addEventListener('mousemove', (e) => {
  mouseX = e.clientX;
  mouseY = e.clientY;
});

function resizeCanvas() {
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
  initWaves();
}

function initWaves() {
  waves = [];
  
  for (let i = 0; i < WAVE_DENSITY; i++) {
    waves.push({
      y: canvas.height * (i + 1) / (WAVE_DENSITY + 1),
      amplitude: WAVE_AMPLITUDE,
      speed: WAVE_SPEED,
      offset: Math.random() * Math.PI * 2
    });
  }
}

function drawWave(wave, time) {
  ctx.beginPath();
  
  const baseY = wave.y;
  const amplitude = wave.amplitude;
  const mouseInfluence = Math.max(0, 1 - Math.hypot(mouseX - canvas.width/2, mouseY - baseY) / 300);
  const currentAmplitude = amplitude * (1 + mouseInfluence);
  
  for (let x = 0; x <= canvas.width; x += 5) {
    const dx = x / canvas.width;
    const offsetY = Math.sin(dx * 10 + time * wave.speed + wave.offset) * currentAmplitude;
    
    if (x === 0) {
      ctx.moveTo(x, baseY + offsetY);
    } else {
      ctx.lineTo(x, baseY + offsetY);
    }
  }
  
  // Compléter la forme jusqu'au bas de l'écran
  ctx.lineTo(canvas.width, canvas.height);
  ctx.lineTo(0, canvas.height);
  ctx.closePath();
  
  // Créer un dégradé pour chaque vague
  const gradient = ctx.createLinearGradient(0, baseY - amplitude, 0, baseY + amplitude);
  
  // Convertir la couleur en RGB pour ajouter la transparence
  const rgbValues = hexToRgb(WAVE_COLOR);
  gradient.addColorStop(0, `rgba(${rgbValues.r}, ${rgbValues.g}, ${rgbValues.b}, 0.7)`);
  gradient.addColorStop(1, `rgba(${rgbValues.r}, ${rgbValues.g}, ${rgbValues.b}, 0.2)`);
  
  ctx.fillStyle = gradient;
  ctx.fill();
}

function hexToRgb(hex) {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result ? {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16)
  } : { r: 0, g: 168, b: 255 };
}

function animate() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  
  const time = Date.now();
  
  // Dessiner chaque vague
  waves.forEach(wave => drawWave(wave, time));
  
  requestAnimationFrame(animate);
}

window.addEventListener('resize', resizeCanvas);

// Initialiser le canvas et commencer l'animation
resizeCanvas();
animate();