const fs = require("fs")
const path = require("path")
const luaJS = require("lua-to-js")

function transpileLuaFiles() {
    const luaDir = path.join(__dirname, 'lua')
    const jsDir = path.join(__dirname, 'js')

    if (!fs.existsSync(jsDir)) {
        fs.mkdirSync(jsDir)
    }

    fs.readdirSync(luaDir).forEach(file => {
        if (file.endsWith('.lua')) {
            const luaCode = fs.readFileSync(path.join(luaDir, file), 'utf8');
            const jsCode = luaJS.transpile(luaCode);

            const jsFileName = file.replace('.lua', '.js');
            fs.writeFileSync(path.join(jsDir, jsFileName), jsCode);
        }
    });

    console.log("Lua files have been eaten by the Balaurus!");
}

transpileLuaFiles();
