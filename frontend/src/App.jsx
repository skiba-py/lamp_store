import { ChakraProvider, Box } from '@chakra-ui/react'
import theme from './theme'
import Header from './components/Header'
import Footer from './components/Footer'
import Catalog from './pages/Catalog'
import Product from './pages/Product'
import Cart from './pages/Cart'
import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom'
import AdminLayout from './components/admin/AdminLayout'
import AdminLogin from './components/admin/AdminLogin'
import AdminProducts from './components/admin/AdminProducts'
import AdminOrders from './components/admin/AdminOrders'

const RequireAdminAuth = ({ children }) => {
  const token = localStorage.getItem('adminToken')
  return token ? children : <Navigate to="/admin/login" />
}

function AppRoutes() {
  const location = useLocation();
  const isAdmin = location.pathname.startsWith('/admin');

  if (isAdmin) {
    return (
      <Routes>
        <Route path="/admin/login" element={<AdminLogin />} />
        <Route path="/admin" element={<RequireAdminAuth><AdminLayout /></RequireAdminAuth>}>
          <Route path="products" element={<AdminProducts />} />
          <Route path="orders" element={<AdminOrders />} />
        </Route>
      </Routes>
    );
  }

  return (
    <>
      <Header />
      <Box id="content" maxW="1200px" w="100%" flex={1} display="flex" flexDirection="column" px={4} py={14} mx="auto">
        <Routes>
          <Route path="/" element={<Catalog />} />
          <Route path="/product/:id" element={<Product />} />
          <Route path="/cart" element={<Cart />} />
        </Routes>
      </Box>
      <Footer />
    </>
  );
}

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Box bg="#E2E2E2">
        <Box id="page" w="100%" display="flex" flexDirection="column">
          <Router>
            <AppRoutes />
          </Router>
        </Box>
      </Box>
    </ChakraProvider>
  )
}
