import { Box, Heading } from '@chakra-ui/react';
import ProductList from '../components/ProductList';
import { products } from '../data';

export default function Catalog() {
  // Заглушка для добавления в корзину
  const handleAddToCart = (product) => {
    alert(`Добавлен в корзину: ${product.name}`);
  };

  return (
    <Box id="catalog" w="100%" minH="100vh">
      <Heading as="h1" size="lg" mb={6} w="100%" textAlign="left">Список товаров</Heading>
      <ProductList products={products} onAddToCart={handleAddToCart} />
    </Box>
  );
}
