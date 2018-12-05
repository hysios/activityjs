import { Button, Container } from 'reactstrap';
import { connect } from 'react-redux';
import React from 'react';

import {Heading} from '../utils'
import { addItemToCart, itemChangeQuantity } from '../../../actions';
import Quantity from '../Quantity';
import Subtotal from '../Subtotal';
import Currency from '../../../Components/Currency'
import styled, { css } from 'react-emotion'

const AddToCartBtn = (props) => (<Button {...props} >添加到购物车</Button>)
const BuyNowBtn = (props) => (<Button {...props} color="primary">购买</Button>)
const Price = ({price, effect}) => (
    <p><Currency as="span"  prefix={"价格"} money={price} /><Effect effects={effect || []} /></p>
)

const fieldsTranslations = {
    "Price": "价格",
    "subtotal": "小计",
}

const Effect = styled(({ effects, className }) => {
    console.log(effects)
    return <span className={className}>{effects.map((eff) => localeEffect(eff)).join(',')}</span>
})`
    color: red;
    margin-left: 0.5rem;
`

function localeEffect(eff) {
    return fieldsTranslations[eff.field] + eff.summary + eff.value
}

const ItemTab = ({item, effect, addItemToCart, itemChangeQuantity}) => {
    return <Container>
        <Heading title={item.title} />
        <img src={item.defaultImage} alt="Image" />
        <Price price={item.price} effect={effect["price"]} />
        <Quantity quantity={item.quantity ? item.quantity : 1} />
        <Subtotal />
        <AddToCartBtn onClick={() => { addItemToCart(item, item.quantity); itemChangeQuantity(1) }} />
        <BuyNowBtn />
    </Container>
}

const mapStateToProps = state => {
    let {effect} = state

    return {effect}
}

export default connect(mapStateToProps, {
    addItemToCart,
    itemChangeQuantity,
})(ItemTab)