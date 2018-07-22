const div = document.createElement('div')
div.innerHTML = 'Hello!'
document.body.append(div)

;(async () => {
    const response = await fetch('http://localhost:3000/polygon_tmp')
    const data = await response.json()
    console.log(data)
})()
