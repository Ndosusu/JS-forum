export function createRegistrationForm(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const title = document.createElement("h2");
    title.textContent = "Inscription / Connexion";
    container.appendChild(title);


        // Formulaire de connexion
    const loginForm = document.createElement("form");
    loginForm.id = "login-form";

    loginForm.innerHTML = `
        <input type="email" id="login-email" placeholder="Email" required>
        <input type="password" id="login-password" placeholder="Password" required>
        <button type="submit" id="login-button">Se connecter</button>
    `;

    container.appendChild(loginForm);


    // Formulaire d'inscription
    const registerForm = document.createElement("form");
    registerForm.id = "register-form";

    registerForm.innerHTML = ` 
        <input type="text" id="username" placeholder="Username" required>
        <input type="email" id="email" placeholder="Email" required>
        <input type="password" id="password" placeholder="Password" required>
        <input type="text" id="firstName" placeholder="First Name" required>
        <input type="text" id="lastName" placeholder="Last Name" required>
        <input type="date" id="age" placeholder="Age" required>
          <label for="gender">Gender:</label>
        <select name="gender">
        <option value="">Please select one…</option>
        <option value="female">Female</option>
        <option value="male">Male</option>
        <option value="non-binary">Non-Binary</option>
        <option value="other">Other</option>
        <option value="Prefer not to answer">Perfer not to Answer</option>
        </select>
        <button type="submit">S'inscrire</button>
    `;

    container.appendChild(registerForm);


    // Bouton pour aller à l'accueil
    const goToIndexButton = document.createElement("button");
    goToIndexButton.textContent = "debug home";
    goToIndexButton.addEventListener("click", () => {
        loadPage("index");
    });

    container.appendChild(goToIndexButton);

    // Ajoute tout au conteneur parent
    parentElement.appendChild(container);


}
