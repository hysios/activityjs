
function pp(ctx) {
    console.table(ctx.item)
    console.table(ctx.order.items)
}

if (window) {
    window.pp = pp
}

export default pp