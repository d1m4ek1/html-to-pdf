const dataFile = {
    fileZIP: new FormData()
}

const inputUploadFile = document.getElementById("input_upload")
const labelUpload = document.getElementById("upload_label")


const buttonSubmit = document.getElementById("submit_upload")

const getFileFromInput = (event) => {
    const file = event.target.files[0]

    if (file.type !== "application/zip") {
        return
    }

    if (file.size > 2147483648) {
        return
    }

    dataFile.fileZIP.set("zip", file)
}

const uploadFile = async () => {
    labelUpload.innerText = "Uploading..."

    const response = await fetch("/api/upload", {
        method: "POST",
        body: dataFile.fileZIP
    })

    const json = await response.json()

    if (json.successfully) {
        labelUpload.innerText = "Uploaded"
        setTimeout(() => {
            labelUpload.innerText = "Upload"
        }, 4000)
    }
}

inputUploadFile.addEventListener("change", getFileFromInput)
buttonSubmit.addEventListener("click", uploadFile)
