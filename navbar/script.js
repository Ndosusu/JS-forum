const list = document.querySelectorAll('.list');
const indicator = document.querySelector('.indicator');

function active() {
    list.forEach((item) => item.classList.remove('active'));
    this.classList.add('active');
    const index = Array.from(list).indexOf(this);
    indicator.style.transform = `translateY(${70 * index}px)`; // Déplacement horizontal
}



list.forEach((item) => item.addEventListener('click', active));

document.querySelectorAll('.navigation .list a').forEach(link => {
    link.addEventListener('click', (event) => {
        event.preventDefault(); // Empêche le comportement par défaut (rechargement de la page)
        const page = event.target.closest('a').getAttribute('href').substring(1); // Récupère l'attribut href sans "#"
        loadPage(page);
    });
});


