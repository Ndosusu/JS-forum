export function createCategories() {
    const cat = document.createElement("div");
    cat.classList.add("categories");
    cat.innerHTML = `
        <ul>
            <li class="category" data-category="bleu">
                <a href="#">
                    <span class="text">Bleu</span>
                </a>
            </li>
            <li class="category" data-category="rouge">
                <a href="#">
                    <span class="text">Rouge</span>
                </a>
            </li>
            <li class="category" data-category="vert">
                <a href="#">
                    <span class="text">Vert</span>
                </a>
            </li>
            <li class="category" data-category="gris">
                <a href="#">
                    <span class="text">Gris</span>
                </a>
            </li>
        </ul>
        <div class="indicator hidden"></div>
    `;

    const container = document.querySelector(".container"); 
    if (container) {
        container.appendChild(cat); // On ajoute la div dans .container
        
        // Ajouter des écouteurs d'événements aux catégories
        const categories = cat.querySelectorAll('.category');
        const indicator = cat.querySelector('.indicator');
        let activeCategory = null;
        
        categories.forEach((category, index) => {
            category.addEventListener('click', () => {
                // Si on clique sur la catégorie déjà active, on la désactive
                if (category === activeCategory) {
                    // Désactivation
                    category.classList.remove('active');
                    indicator.classList.add('hidden');
                    activeCategory = null;
                    
                    // Réinitialiser le filtrage après un court délai pour l'animation
                    setTimeout(() => {
                        filterPostsByCategory('all');
                    }, 300); // Délai correspondant à l'animation
                } else {
                    // Rétablir la transition pour assurer l'animation entre catégories
                    indicator.style.transition = 'transform 0.5s';
                    
                    // Positionner l'indicateur à l'emplacement correct avec animation
                    indicator.style.transform = `translateY(calc(50px * ${index}))`;
                    
                    // Supprimer la classe 'active' de toutes les catégories
                    categories.forEach(c => c.classList.remove('active'));
                    
                    // Ajouter la classe 'active' à la catégorie cliquée
                    category.classList.add('active');
                    
                    // Afficher l'indicateur s'il était caché
                    indicator.classList.remove('hidden');
                    
                    // Mettre à jour la catégorie active
                    activeCategory = category;
                    
                    // Récupérer la valeur de la catégorie
                    const categoryValue = category.getAttribute('data-category');
                    
                    // Filtrer les posts
                    filterPostsByCategory(categoryValue);
                }
            });
        });
    }
}