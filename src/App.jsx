import { ChakraProvider, Box } from '@chakra-ui/react'
import theme from './theme'
import Header from './components/Header'
import Footer from './components/Footer'
import Catalog from './pages/Catalog'
import Product from './pages/Product'
import Cart from './pages/Cart'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Box bg="#E2E2E2">
        <Box
          id="page"
          maxW="1200px"
          w="100%"
          mx="auto"
          px={4}
          display="flex"
          flexDirection="column"
        >
          <Router>
            <Header />
            <Box flex={1} display="flex" flexDirection="column">
              <Routes>
                <Route path="/" element={<Catalog />} />
                <Route path="/product/:id" element={<Product />} />
                <Route path="/cart" element={<Cart />} />
              </Routes>
            </Box>
            <Footer />
          </Router>
        </Box>
      </Box>
    </ChakraProvider>
  )
}
