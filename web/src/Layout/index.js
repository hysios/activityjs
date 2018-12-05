import React, {Component} from 'react'
import styled, { css } from 'react-emotion'
import Page from '../Components/Page';
import Code from '../Components/Code';
import SplitHandler from '../Components/SplitHandler';

const style=css`
  display: flex;
  flex-wrap: nowrap;
`

const Layout = (props) => (
    <div className={style} >
        <SplitHandler>

            <Page >{props.children}</Page>
            <Code />
        </SplitHandler>
    </div>
)

export default Layout