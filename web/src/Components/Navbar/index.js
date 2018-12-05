import {
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    NavLink,
    UncontrolledDropdown,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';
import { connect } from 'react-redux'
import React, { Component } from 'react';

import MyCartButton from '../MyCart/Button';
import SearchProducts from '../SearchProducts';
import { fetchExample, setCode } from '../../actions/code'
import { source } from '../Code';

class MainNavbar extends Component {
    state = {
        isOpen: false
    }

    handleActivityType(filename) {
        this.props.fetchExample(filename)
    }

    render() {
        return <Navbar color="light" light expand="md">
            <NavbarBrand href="/">ActivityJS 活动 Demo</NavbarBrand>
            <NavbarToggler onClick={this.toggle} />
            <Collapse isOpen={this.state.isOpen} navbar>
                <SearchProducts />

                <Nav className="ml-auto" navbar>
                    <NavItem>
                        <NavLink href="/components/">Components</NavLink>
                    </NavItem>
                    <NavItem>
                        <NavLink href="https://github.com/reactstrap/reactstrap">GitHub</NavLink>
                    </NavItem>
                    <MyCartButton />
                    <UncontrolledDropdown nav inNavbar>
                        <DropdownToggle nav caret>
                            活动类型
                        </DropdownToggle>
                        <DropdownMenu right>
                            <DropdownItem onClick={() => this.handleActivityType("standard.go")}>
                                Standard
                            </DropdownItem>
                            <DropdownItem onClick={() => this.handleActivityType("for_test.go")}>
                                For Debug
                            </DropdownItem>
                            <DropdownItem divider />
                            <DropdownItem onClick={() => this.props.setCode(source) }>
                                Reset
                            </DropdownItem>
                        </DropdownMenu>
                    </UncontrolledDropdown>
                </Nav>
            </Collapse>
        </Navbar>
    }
}

export default connect(null, {
    fetchExample,
    setCode,
})(MainNavbar)