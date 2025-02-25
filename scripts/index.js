export function createIndex(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const title = document.createElement("h2");
    title.textContent = "Accueil";
    container.appendChild(title);

    // Bouton pour retourner à la page d'inscription
    const backButton = document.createElement("button");
    backButton.textContent = "Retour à l'inscription";
    backButton.addEventListener("click", () => {
        loadPage("registration");
    });

    const profileButton = document.createElement("button");
    profileButton.textContent = "Votre profile";
    profileButton.addEventListener("click", () => {
        loadPage("profile");
    });

    container.appendChild(backButton);
    container.appendChild(profileButton);

    // New post
    const newPost = document.createElement("div");
    newPost.classList.add("newpostcontainer");
    container.appendChild(newPost);

    const newPostTitle = document.createElement("input");
    newPost.appendChild(newPostTitle);

    const newPostContent = document.createElement("textarea");
    newPost.appendChild(newPostContent);

    const newPostButton = document.createElement("button");
    newPostButton.textContent = "Post";
    newPost.appendChild(newPostButton);

    // Posts
    const posts = document.createElement("div");
    posts.classList.add("posts");
    container.appendChild(posts);

    const post = document.createElement("div");
    post.classList.add("post");
    posts.appendChild(post);

    const postTitle = document.createElement("h3");
    postTitle.textContent = "title";
    post.appendChild(postTitle);

    const postBody = document.createElement("div");
    postBody.textContent = "testPost";
    post.appendChild(postBody);

    const answers = document.createElement("div");
    answers.textContent = "testanswers";
    post.appendChild(answers);

    const answerIt = document.createElement("textarea");
    post.appendChild(answerIt);

    const answerButton = document.createElement("button");
    answerButton.textContent = "send";
    post.appendChild(answerButton);

    // Ajoute tout au conteneur parent (et non document.body)
    parentElement.appendChild(container);
}


