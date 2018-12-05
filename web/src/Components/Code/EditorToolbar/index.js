import { ButtonGroup, ButtonToolbar, Button } from 'reactstrap';
import {connect} from 'react-redux'
import { css } from 'react-emotion';
import React from 'react';

import { addItemToCode, addProductToCode } from '../../../actions/code';

const toolbarStyle = css`
    margin: 0.5rem;
    margin-left: 0;
`

const EditorToolbar = ({item, addProductToCode, addItemToCode}) => (
    <ButtonToolbar className={toolbarStyle} >
        <ButtonGroup size="sm">
            <Button >添加购物车到订单</Button>
            <Button onClick={() => addProductToCode(item)}>添加当前产品</Button>
            <Button onClick={() => addItemToCode({...item, id: item.productId, produceDate: 0})}>添加当前 Item</Button>
            <Button>4</Button>
        </ButtonGroup>
        <ButtonGroup size="sm">
            <Button>5</Button>
            <Button>6</Button>
            <Button>7</Button>
        </ButtonGroup>
        <ButtonGroup size="sm">
            <Button>8</Button>
        </ButtonGroup>
    </ButtonToolbar>

)

const mapStateToProps = state => {
    let {item} = state

    return {item}
}

export default connect(mapStateToProps, {
    addProductToCode,
    addItemToCode
})(EditorToolbar)