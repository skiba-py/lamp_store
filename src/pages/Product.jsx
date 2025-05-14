import { Box, Heading, Text, Image, Button, Flex } from '@chakra-ui/react';
import { useParams, Link } from 'react-router-dom';
import { products } from '../data';

export default function Product({ onAddToCart }) {
  const { id } = useParams();
  const product = products.find(p => String(p.id) === id);

  if (!product) return <Box>Товар не найден</Box>;

  return (
    <Box id="product" minH="100vh" w="100%">
      <Flex gap={10} align="flex-start" w="100%" maxW="1200px">
        <Image src={product.image} alt={product.name} boxSize="320px" objectFit="cover" borderRadius="md" />
        <Box flex={1}>
          <Heading as="h2" size="lg">{product.name}</Heading>
          <br/>
          <Flex align="center" gap={4} mb={2}>
            <Text fontSize="xl" fontWeight="bold">{product.price} ₽</Text>
            <Button onClick={() => onAddToCart(product)} variant="outline" bg="white">Добавить в корзину</Button>
          </Flex>
          <Text fontWeight="bold" mb={2}>Описание</Text>
          <Text color="gray.600">{product.description}</Text>
          <Link to="/" style={{ marginTop: 16, display: 'inline-block', color: '#3182ce' }}>← Назад в каталог</Link>
        </Box>
      </Flex>
    </Box>
  );
}
