import { generateCode, generateInSliceCode } from "../../utils/compile"
import { capitalize } from "lodash"

const reReact = /^\/\/go:react:(\w+)$/gm
const reImportGroup = /import\s+(\((.|\n)*?\))/gm
const rePackage = /^package \w+$/gm
const reImport = function (str) {
    let re = [reImportGroup, /^import\s+("|\().*?\1/gm]
    let res = [ re[0].exec(str), re[1].exec(str)]
    return res[0] ? [res[0], re[0]] : [res[1],re[1]]
}

const reScopeHeader = /^var (\w+) = (\w+.)?(\w+){$/g

const zero = { start: 0, end: 0 }

export default class Compiler {

    static parse(source) {
        return new Compiler(source)
    }

    constructor(source) {
        this.reReact = /^\/\/go:react:(\w+)$/gm
        this.rePackage = /^package \w+$/gm


        this.source = source
        this.blocks = {
            REACT: {
                products: zero,
                product: zero,
                order: zero,
                item: zero
            },
            IMPORT_BLOCKS: zero,
            IMPORTS: []
        }


        let matchs = this.parseReactTag(this.source)
        this.parseImportHeader(this.source)
        this.parseImports(this.source)
    }

    parseReactTag(source) {
        let matchs = []
        let match = this.reReact.exec(source)
        let blocks = this.blocks.REACT
        while (match != null) {
            blocks[match[1]] = {
                start: match.index,
                end: this.reReact.lastIndex + 1
            }
            switch(match[1]) {
            case "products": {
                let [_, pos] = this.parseScope(source, match.index)
                blocks['products']['insert'] = pos
            }
            }

            match = this.reReact.exec(source);
        }
        return matchs
    }


    parseBracket(content, brackets = '{}', fn) {
        let stack = []
        for (let i = 0 ;i < content.length; i++) {
            let c = content[i]
            if (c == brackets[0]) {
                stack.push(i)
            } else if (c == brackets[1] && !!stack.length) {
                let start = stack.pop()
                fn(stack.length, content.slice(start+1,i), i)
                if (stack.length === 0) break
            }
        }
    }

    parseScope(source, start) {
        let code = source.slice(start)
        let reslt = []

        this.parseBracket(code, '{}', (i, c, pos) => reslt.push([i, c, pos]))
        let [_, c , lastPos] = reslt[reslt.length - 1]
        return [c, start + lastPos]
    }

    parseImportHeader(source) {
        let [match, re] = reImport(source)
        if (!match) {
            match = this.rePackage.exec(source)
            if (match == null) {
                // throw Error("不能够解析到 package")
                console.error(new Error("不能够解析到 package"))
                return
            }
            this.blocks.IMPORT_BLOCKS = {
                created: true,
                start: this.rePackage.lastIndex + 1,
                end: this.rePackage.lastIndex + 1,
            }
        } else {
            this.blocks.IMPORT_BLOCKS = {
                start: match.index,
                end: re.lastIndex + 1,
                content: {
                    start: re.lastIndex - match[1].length,
                    end: re.lastIndex + 1
                }
            }
        }
    }

    parseImports (source) {
        let imports = "", block = this.blocks.IMPORT_BLOCKS

        if (block.created) {
            return
        }

        imports = source.slice(block.content.start, block.content.end)

        if (imports[0] === '"') {
            this.blocks.IMPORTS = [{
                id: imports.slice(1, imports.length-2),
            }]
        } else if (imports[0] === "(") {
            imports = imports.slice(1, imports.length-2)
            let imps = imports.split("\n")
            this.blocks.IMPORTS = imps.filter(m => this.pure(m)).map(m => ({id: this.pure(m)}))
        }
    }

    pure(m) {
        m = m.replace(/"/g, '').trim()
        console.log('pure', m)
        return m.replace(/"/g, '').trim()
    }

    gencode(model, val) {
        if (typeof val === 'string') {
            return val
        } else if (typeof val === 'object') {
            return generateCode(model, val)
        }
    }

    genSliceMembercode(model, val) {
        if (typeof val === 'string') {
            return val
        } else if (typeof val === 'object') {
            return generateInSliceCode(model, val)
        }
    }

    eval(context) {

        let react = this.blocks.REACT
        let offset = 0
        let newCode = ""
        let { imports } = context
        // let addGroup = false

        // let pos = this.blocks.IMPORT_BLOCKS
        // let content = {}
        // console.log(this.blocks)
        // if (pos.created) {
        //     let before = this.source.slice(offset, pos.end)
        //     offset += pos.end
        //     newCode += before + "\nimport (\n"
        //     addGroup = true
        //     content = {
        //         start: offset,
        //         end: offset,
        //     }
        // } else {
        //     content = this.blocks.IMPORT_BLOCKS.content
        // }

        // Object.keys(imports || {}).forEach(key => {
        //     let val = context[key]
        //     if (val) {
        //         let imp = imports[key]

        //         if (!this.hasImport(imp)) {
        //             let before = this.source.slice(content.offset, content.end)
        //             offset = pos.end
        //             newCode += before + "\t" + imp
        //         }

        //     }
        // })

        // if (addGroup) {
        //     offset += 1
        //     newCode += "\n)\n"
        // }

        Object.keys(react).forEach(key => {
            let pos = react[key]
            let val = context[key]

            switch(key){
            case "product": {
                pos = react['products']
                if (this.hasReactBlock('products')) {
                    let before = this.source.slice(offset, pos.insert)
                    offset = pos.insert
                    newCode += before +  this.genSliceMembercode(capitalize(key), val) + "\n"
                } else {
                    let before = this.source.slice(offset, pos.end)
                    offset = pos.end
                    newCode += before + this.gencode(capitalize(key), val)
                }
                break
            }
            default: {
                if (pos.start > 0 && val) {
                    let before = this.source.slice(offset, pos.end)
                    offset = pos.end
                    newCode += before + this.gencode(capitalize(key), val)
                }
            }
            }
        });

        return newCode + this.source.slice(offset)
    }

    hasReactBlock(name) {
        let blocks = this.blocks.REACT
        if (blocks[name] && blocks[name].start > 0) {
            return true
        }
        return false
    }

    evalImports(context) {

    }

    hasImport(id) {
        let imports = this.blocks.IMPORTS
        for (let i in imports) {
            if (imports[i].id == id) {
                return true
            }
        }
        return false
    }

}
