import './index.css'

import { connect } from 'react-redux';
import React, { Component } from 'react'
import SplitPane from 'react-split-pane'
import { splitCode } from '../../actions/code';
import styled, { css } from 'react-emotion'

const SplitContainer = styled('div')`
    padding-left: 8px;
    &:hover {
        background: #ccc;
    }
`

const Handler = styled('div')`
    height: 100%;
    width: 8px;
`


class SplitHandler extends Component {
    state = {}

    handleChange = (size) => {
        console.log(size)
        this.props.splitCode(window.innerWidth - size)
    }

    render() {
        return <SplitPane split="vertical" defaultSize="70%" onChange={this.handleChange} >{this.props.children}</SplitPane>
    }
}



export default connect(null, { splitCode })(SplitHandler)