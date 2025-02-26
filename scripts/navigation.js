export function createNavigation() {
    const nav = document.createElement("div");
    nav.classList.add("navigation");
    nav.innerHTML = `
        <ul>
            <li class="list active">
                <a href="#index">
                    <span class="icon">
                        <ion-icon name="home-outline"></ion-icon>
                    </span>
                    <span class="text">Home</span>
                </a>
            </li>
            <li class="list">
                <a href="#">
                    <span class="icon">
                        <ion-icon name="book-outline"></ion-icon>
                    </span>
                    <span class="text">About</span>
                </a>
            </li>
            <li class="list">
                <a href="#contact">
                    <span class="icon">
                        <ion-icon name="chatbox-ellipses-outline"></ion-icon>
                    </span>
                    <span class="text">Contact</span>
                </a>
            </li>
            <li class="list">
                <a href="#">
                    <span class="icon">
                        <ion-icon name="settings-outline"></ion-icon>
                    </span>
                    <span class="text">Setting</span>
                </a>
            </li>
            <li class="list">
                <a href="#profile">
                    <span class="icon">
                        <ion-icon name="person-outline"></ion-icon>
                    </span>
                    <span class="text">Profile</span>
                </a>
            </li>
            <div class="indicator"></div>
        </ul>
    `;

    document.body.prepend(nav);

    // Charger dynamiquement le script de la navbar après l'insertion de la navigation
    const script = document.createElement("script");
    script.src = "navbar/script.js";
    script.type = "module";
    document.body.appendChild(script);
}
