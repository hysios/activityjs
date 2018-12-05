import {
    Button,
    Card,
    CardText,
    CardTitle,
    Col,
    Container,
    Nav,
    NavItem,
    NavLink,
    Row,
    TabContent,
    TabPane
} from 'reactstrap';
import { connect } from 'react-redux'
import React, { Component } from 'react'
import styled, { css, cx } from 'react-emotion'
import { sprintf } from 'sprintf-js'
import { productToItem} from '../../actions'
import { genItemCode, genProductCode, genUserCode } from '../../actions/code'
import Quantity from './Quantity'
import Subtotal from './Subtotal'
import { Currency } from './utils'
import ItemTab from './Tabs/ItemTab'
import { Heading } from './utils'

class Tabs extends Component {
    state = {
        activeTab: '1',
    }

    toggle = (tab) => {
        if (this.state.activeTab !== tab) {
            this.setState({
                activeTab: tab
            });
        }
    }

    render() {
        return <div >
            <Nav tabs>
                <NavItem>
                    <NavLink
                        className={cx({ active: this.state.activeTab === '1' })}
                        onClick={() => { this.toggle('1'); }}
                    >
                        产品
                </NavLink>
                </NavItem>
                <NavItem>
                    <NavLink
                        className={cx({ active: this.state.activeTab === '2' })}
                        onClick={() => { this.toggle('2'); }}
                    >
                        购物车
                </NavLink>
                </NavItem>
                <NavItem>
                    <NavLink
                        className={cx({ active: this.state.activeTab === '3' })}
                        onClick={() => { this.toggle('3'); }}
                    >
                        订单
                </NavLink>
                </NavItem>
            </Nav>
            <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="1">
                    <ItemTab item={this.props.item} />
                </TabPane>
                <TabPane tabId="2">
                    <Row>
                        <Col sm="6">
                            <Card body>
                                <CardTitle>Special Title Treatment</CardTitle>
                                <CardText>With supporting text below as a natural lead-in to additional content.</CardText>
                                <Button>Go somewhere</Button>
                            </Card>
                        </Col>
                        <Col sm="6">
                            <Card body>
                                <CardTitle>Special Title Treatment</CardTitle>
                                <CardText>With supporting text below as a natural lead-in to additional content.</CardText>
                                <Button>Go somewhere</Button>
                            </Card>
                        </Col>
                    </Row>
                </TabPane>
            </TabContent>
        </div>
    }
}

class Content extends Component {
    render() {
        let { item } = this.props;
        return <Container>
            <Tabs item={item} />
        </Container>
    }
}

const mapStateToProps = state => {
    let { item } = state
    return { item }
}

export default connect(mapStateToProps)(Content)