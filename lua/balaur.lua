#return "<h1>Balaurus!</h1>"

local Button = require("components/button")

local balaurButton = Button.new("Balaurus", "handleClick()")

function handleClick()
    print("Balaurus: Growl!!!")

    document.getElementById("message").innerHTML = "Balaurus!!"
end

_G.components = {
    balaurButton = balaurButton
}
