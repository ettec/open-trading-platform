import OrderBlotter from "../components/OrderBlotter"
import { RootState, setSelectedOrder } from "../redux"
import { connect } from 'react-redux'

const mapStateToProps = (state: RootState) => ({
    selectedOrder : state.selectedOrder
  })

const mapDispatchToProps = {
    onOrderSelected: setSelectedOrder
}

export default connect(mapStateToProps,  mapDispatchToProps)(OrderBlotter)