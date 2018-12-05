import {
  Form,
  FormGroup,
  Input,
  ListGroup,
  ListGroupItem,
  Popover,
  PopoverBody,
} from 'reactstrap';
import { Redirect } from 'react-router'
import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import { search } from '../../search'
import styled, { css } from 'react-emotion'
import Currency from '../Currency'

const highlight=css`
    display: flex!important;
    width: 320px;
    em {
        background: rgba(255, 193, 7, 0.65);
        color: #0062cc;
    }
`

const Icon = styled('img')`
    width: 1.5rem;
    height: 1.5rem;
    margin-right: 0.5rem;
`

const itemStyle = css`
    overflow: hidden;
    min-width: 0px;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
`

const moneyStyle = css`
    flex: 1;
    text-align: right;
    width: auto;
`

const pullStyle = css`
    max-width: 375px!important;
`

const Title = ({item}) => (
    <span className={itemStyle} dangerouslySetInnerHTML={{__html: item._highlightResult.name.value}} />
)


const SearchItem = ({item, index, select}) => (
    <ListGroupItem tag="a" href={`/products/${item.id}`} key={item.id} className={highlight} active={index == select} >
        <Icon src={item.imageUrls[0]}  /><Title item={item} /><Currency as='span' className={moneyStyle} money={item.price} />
    </ListGroupItem>
)

class SearchProducts extends Component {
    state = { popoverOpen: false, items: [], select: 0 }

    toggle = () => {
        this.setState({
            popoverOpen: !this.state.popoverOpen
        });
    }

    handleSearch = async (e) => {
        let result = this.getResults(await search(e.target.value))
        let state = {popoverOpen: true}
        if (result.length > 0) {
            state["items"] = result
        }
        this.setState(state)
    }

    handleKeySelect = (e) => {
        switch (e.which) {
        case 40: {// down {
            let {items, select} = this.state
            select = select < items.length ? select+1 : 0
            this.setState({
                select,
                popoverOpen: true,
            })
            break
        }
        case 38: { //up
            let {items, select} = this.state
            select = select > 0 ? select-1 : 0
            this.setState({
                select
            })
            break
        }
        case 13: { // enter
            let {items, select} = this.state
            let {history} = this.props
            history.push(`/products/${items[select].id}`, null)
            this.setState({
                popoverOpen: false
            })
        }
        default:
            console.log(e.which)
        }
    }

    getResults(res) {
        let hits = res.hits || []

        return hits
    }

    render() {
        let { items, select, redirect } = this.state

        return <Form className="form-inline" onSubmit={(e) => e.preventDefault() }>
            {redirect ? <Redirect push to={redirect} /> : null}
            <FormGroup>
                <Input type="search"
                    name="search"
                    id="product-search"
                    placeholder="搜索产品"
                    onChange={this.handleSearch}
                    onKeyDown={this.handleKeySelect} />
            </FormGroup>
            <Popover className={pullStyle} placement="bottom" isOpen={this.state.popoverOpen} target="product-search" toggle={this.toggle}>
                <PopoverBody>
                    <ListGroup >
                        {items.map((item, i) => <SearchItem key={item.id} item={item} index={i} select={select} />)}
                    </ListGroup>
                </PopoverBody>
            </Popover>
        </Form>
    }
}

export default withRouter(SearchProducts)