import { Badge } from 'reactstrap';
import { connect } from 'react-redux'
import React from 'react';
import styled, { css } from 'react-emotion'
import Humanize from 'humanize-plus'
import Transition from 'react-transition-group/Transition';

const mapCompileStateToProps = state => {
    let { compile, api: {working} } = state

    return {...compile, working}
}

const duration = 300;

const defaultStyle = {
    transition: `opacity ${duration}ms ease-in-out`,
    opacity: 0,
}

const transitionStyles = {
    entering: { opacity: 0 },
    entered:  { opacity: 1 },
};

const TransitionStatus = ({show, color, text, errors}) => (
    <Transition in={show} timeout={{enter: 0}}
                unmountOnExit >
        {(state) => {
            console.log(state)
            return (<React.Fragment>
                <StatuText style={{...defaultStyle, ...transitionStyles[state]}} color={color} text={text} />
                { errors ? <ErrorsList errors={errors} /> : null }
            </React.Fragment>)
        }}
    </Transition>
)

const Status = ({ state, size, duration, errors = null, working }) => {
    let text, show = false, color;
    console.log(state, working)
    if (working.COMPILE_GOCODE) {
        show = false
    } else {
        switch (state) {
        case "success":
            text = "编译成功 " + duration + "ms" + " 共 " + Humanize.fileSize(size)
            color = "success"
            show = true
            break
        case "error":
            text = "编译失败"
            color = "danger"
            show = true
            break
        case null:
        default:
            show = false
        }
    }
    return <TransitionStatus show={show} color={color} text={text} errors={errors} />
}

const ErrorsList = (({errors = []}) => {
    return errors.map((err) => (<Badge color="light" pill>{`${basename(err.filename)}:${err.line} ${err.msg}`}</Badge>))
})

function basename(file) {
    let pos = file.lastIndexOf("/")
    if (pos > 0) {
        return file.slice(pos+1)
    }
    return file
}


const StatuText = styled(({style,  color, text, className }) => {
    return <Badge style={style} className={className} color={color}>{text}</Badge>
})`
    margin-left: 1rem;
    label: status-text;
`
const CompileStatus = connect(mapCompileStateToProps)(Status)

export default CompileStatus
