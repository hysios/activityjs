import {
  Button,
  ButtonDropdown,
  DropdownItem,
  DropdownMenu,
  DropdownToggle,
} from 'reactstrap';
import { connect } from 'react-redux';
import React, { Component } from 'react'

class CompileButton extends Component {
    state = { dropdownOpen: false }

    toggle = () => {
        this.setState({
            dropdownOpen: !this.state.dropdownOpen
        });
    }

    render() {
        let {working, onClick} = this.props

        return (
            <ButtonDropdown size="lg" disabled={working} isOpen={this.state.dropdownOpen} toggle={this.toggle}>
                <Button onClick={onClick} disabled={working} color="success">编&nbsp;&nbsp;&nbsp;&nbsp;译</Button>
                <DropdownToggle disabled={working} caret color="success"  />
                <DropdownMenu right>
                    <DropdownItem disabled={working} onClick={(e) => onClick(e, true)} >编&nbsp;&nbsp;&nbsp;&nbsp;译（调试）</DropdownItem>
                </DropdownMenu>
            </ButtonDropdown>
        )
    }
}

const mapStateApiToProps = state => {
    let {api} = state
    return {working: api.working["COMPILE_GOCODE"]}
}

export default connect(mapStateApiToProps)(CompileButton)