import { Button, Badge } from 'reactstrap'
import { connect } from 'react-redux'
import Compiler from './Compiler'
import { compile, vet } from '../../actions/compile'
import { modifyCode } from '../../actions/code'

import styled, { css } from 'react-emotion'
import React, { Component } from 'react'
import './index.css'
import EditorToolbar from './EditorToolbar'
import CompileButton from './CompileButton'
import CompileStatus from './CompileStatus'

// import 'codemirror/mode/go/go'
// import 'codemirror/mode/javascript/javascript'
// import 'codemirror/mode/markdown/markdown'
// import 'codemirror/mode/xml/xml'

import brace from 'brace';
import AceEditor from 'react-ace';

import 'brace/mode/golang';
import 'brace/theme/monokai';
import 'brace/theme/github';

const style = css`
    width: 600px;
    .CodeMirror {
        height: calc(100vh - 56px - 120px - 31px)!important;
    }
`

const CodePanel = styled("div")`
    padding: 1rem;
    height: 120px;
    background: #efefef;
`
export const source = `// remote gopherjs generator

package main

//go:react:product
//go:react:item
//go:react:order

func main() {

}
`
function mycode() {
    return localStorage.getItem("mycode")
}

function savecode(code) {
    return localStorage.setItem("mycode", code)
}

class Code extends Component {
    state = {
        mode: 'go'
    }

    code = mycode() || source

    // componentWillReceiveProps(nextProps) {
    //     let codeOrNil, statusOrNil, annotationsOrNil, widthOrNil
    //     let newcode

    //     if (nextProps.code) {
    //         codeOrNil = "code"
    //     }

    //     if (nextProps.status) {
    //         statusOrNil = "status"
    //         // if (nextProps.status == "error") {
    //         //     this.props.vet(nextProps.code || this.state.code)
    //         // }
    //     }

    //     let {product, item} = nextProps

    //     if (product || item) {
    //         let context = {
    //             ...this.state.context,
    //             product,
    //             item
    //         }

    //         newcode = this.buildCode(this.state.code, context)
    //     }

    //     if (nextProps.annotations) {
    //         annotationsOrNil = "annotations"
    //     }

    //     if (nextProps.width) {
    //         widthOrNil = "width"
    //     }

    //     this.setState({
    //         [codeOrNil]: nextProps.code || newcode,
    //         [statusOrNil]: nextProps.status,
    //         [annotationsOrNil]: nextProps.annotations,
    //         [widthOrNil]: nextProps.width,
    //     })

    // }

    buildCode(code, _context) {
        console.log('buildCode')
        let parser = Compiler.parse(code)
        let context = {
            ..._context,
            imports: {
                product: '"activityjs.io/serve/model"'
            }
        }
        // this.updateCode(parser.eval(context))
        return this.updateCode(parser.eval(context))
    }

    updateCode = (newCode) => {
        this.code = newCode
        // this.setState({
        //     code: newCode,
        // })

        this.props.modifyCode()
        savecode(this.code)
        return this.code
    }

    insertCode = () => {
        let {product, item} = this.props
        console.log('insertCode')

        if (product || item) {
            let context = {
                ...this.state.context,
                product,
                item
            }
            return this.buildCode(this.code, context)
        }
        return this.code
    }

    render() {
        var options = {
            lineNumbers: true,
        };

        let { code, annotations, width } = this.props;
        annotations = annotations || []
        let markers =  annotations.map(({row, colume}) => ({startRow: row, startCol: 0, endRow: row, endCol: 1000, type: 'background', className:'error-marker'}) )
        let codeOrNil;
        if (code) {
            codeOrNil = code
            this.updateCode(code)
        } else {
            codeOrNil = this.insertCode(this.code)
        }
        return <div>
            <EditorToolbar />
            <AceEditor
                mode="golang"
                theme="github"
                name="GOLANG_EDITOR"
                value={codeOrNil}
                annotations={annotations}
                markers={markers}
                editorProps={{ $blockScrolling: true }}
                height={"calc(100vh - 56px - 120px - 31px)"}
                width={width || "600px"}
                onChange={this.updateCode}
            />
            <CodePanel>
                <CompileButton onClick={(e, debug) => this.props.compile(codeOrNil, debug)} />
                <CompileStatus />
            </CodePanel>
        </div>
    }
}

const mapStateToProps = state => {
    let { code, compile } = state

    return {
        code: code.code,
        product: code.product,
        item: code.item,
        order: code.OrderCode,
        status: compile.state,
        annotations: compile.annotations,
        width: code.codeEditorWidth,
    }
}

export default connect(mapStateToProps, { compile, modifyCode, vet })(Code)
