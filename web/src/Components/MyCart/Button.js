import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
} from 'reactstrap';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import styled, { css } from 'react-emotion'
import { Currency } from '../../Pages/Product/utils';
import { sprintf } from 'sprintf-js'


const Money = styled('span')`
`

const CartTotal = ({ price }) => (
    isNaN(price) || price == 0 ? null : <Money>￥ {sprintf("%0.2f", price)}</Money>
)

const ItemQuantity = styled('span')`
    border-left: 1px solid #dee2e6!important;
    width: 30px;
    text-align: right;
`

const ItemRow = styled('span')`
    max-width: 280px;
    display: flex;
    flex-wrap: nowrap;
`

const ItemTitle = styled('span')`
    flex: 2;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
`

class MyCartButton extends Component {
    state = {
        dropdownOpen: false
    }

    toggle = () => {
        this.setState({ dropdownOpen: !this.state.dropdownOpen })
    }

    total() {
        let { items } = this.props
        return items.reduce((s, m) => (m.quantity * m.price), [])
    }

    render() {
        let { items } = this.props
        return <div>
            <ButtonDropdown isOpen={this.state.dropdownOpen} toggle={this.toggle}>
                <DropdownToggle caret>
                    购物车 <CartTotal price={this.total()}></CartTotal>
                </DropdownToggle>
                <DropdownMenu right>
                    {items.length > 0 ? <DropdownItem header>购物车的产品</DropdownItem> : <DropdownItem header>购物车是空的</DropdownItem>}
                    {items.map((item) => (
                        <DropdownItem key={item.productId}><ItemRow > <ItemTitle>{item.title} </ItemTitle><ItemQuantity>{item.quantity}</ItemQuantity> </ItemRow> </DropdownItem>
                    ))}
                </DropdownMenu>
            </ButtonDropdown>
        </div>
    }
}

const mapStateToProps = state => {
    let { cart } = state

    return { items: cart.items }
}

export default connect(mapStateToProps)(MyCartButton)