import React, {Component} from 'react'
import styled, { css } from 'react-emotion'

const style=css`
    flex: 1
`

const Page = ({children}) => (
    <div className={style}>
        {React.Children.toArray(children)}
    </div>
)

export default Page;