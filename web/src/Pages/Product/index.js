import { connect } from 'react-redux';
import { withRouter } from 'react-router'
import React, {Component} from 'react';

import Content from './Content';
import Layout from '../../Layout/index';
import {fetchProduct} from '../../actions'

class Product extends Component {
    state = {
        productId: this.props.match.params.id
    }

    componentDidMount () {
        let { match: {
            params: {
                id
            }
        } } = this.props
        this.oldProductId = id
        this.props.fetchProduct(id)
    }

    static getDerivedStateFromProps(props, state) {
        if (props.match.params.id !=  state.productId)  {

            return {
                productId: props.match.params.id,
            }
        }
        return null
    }

    render() {
        let { match: {
            params: {
                id
            }
        } } = this.props
        if (this.oldProductId != id ) {
            this.oldProductId = id
            this.props.fetchProduct(id)
        }

        return (<Layout>
            <Content />
        </Layout>)
    }
}

export default withRouter(connect(null, {fetchProduct})(Product))
