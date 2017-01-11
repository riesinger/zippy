document.addEventListener("DOMContentLoaded", function() {
    var simplemde = new SimpleMDE({
                element: document.getElementById("editor"),
                autoDownloadFontAwesome: false,
                toolbar: false
            });

    var saveButton = document.body.querySelector("button#save");
    var titleInput = document.body.querySelector("input#post-title");
    var authorInput = document.body.querySelector("input#author-id");
    var tmpl = document.body.querySelector("input#template");
    var path = document.body.querySelector("input#path");
    var draftCheckbox = document.body.querySelector("input#draft-checkbox");

    saveButton.addEventListener("click", function() {
        var postContent = simplemde.value();
        var d = new Date();

        var postData = JSON.stringify({
            title: titleInput.value,
            isDraft: draftCheckbox.checked,
            authorID: authorInput.value,
            template: tmpl.value,
            path: path.value,
            markdownBody: postContent
        });

        console.log(postData);

        fetch("/api/article/new", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: postData
        })
            .then(function(response) {
                return response.json();
            })
            .then(function(response) {

            })
            .catch(function(err) {
                console.error(err);
            });
    });

});
