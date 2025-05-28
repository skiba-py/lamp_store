import { Box, Heading, Text, Image, Button, Flex, Spinner, useToast } from '@chakra-ui/react';
import { useParams, Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { productsApi, ordersApi } from '../services/api';

export default function Product() {
  const { id } = useParams();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [adding, setAdding] = useState(false);
  const toast = useToast();

  useEffect(() => {
    const fetchProduct = async () => {
      setLoading(true);
      try {
        const data = await productsApi.getProductById(id);
        setProduct(data);
      } catch (err) {
        setError('Товар не найден');
      } finally {
        setLoading(false);
      }
    };
    fetchProduct();
  }, [id]);

  const handleAddToCart = async () => {
    setAdding(true);
    try {
      await productsApi.checkProductAvailability(product.id, 1);
      const orderData = {
        user_id: "123e4567-e89b-12d3-a456-426614174000",
        status: "pending",
        total: product.price,
        items: [
          {
            product_id: product.id.toString(),
            quantity: 1,
            price: product.price
          }
        ]
      };
      await ordersApi.createOrder(orderData);
      toast({
        title: "Товар добавлен в корзину",
        description: `${product.name} успешно добавлен в корзину`,
        status: "success",
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: "Ошибка",
        description: error.message || "Не удалось добавить товар в корзину",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setAdding(false);
    }
  };

  if (loading) return <Box textAlign="center" py={10}><Spinner size="xl" /></Box>;
  if (error) return <Box textAlign="center" py={10} color="red.500">{error}</Box>;

  return (
    <Box id="product" minH="100vh" w="100%">
      <Flex gap={10} align="flex-start" w="100%" maxW="1200px">
        <Image src={product.image} alt={product.name} boxSize="320px" objectFit="cover" borderRadius="md" />
        <Box flex={1}>
          <Heading as="h2" size="lg">{product.name}</Heading>
          <br/>
          <Flex align="center" gap={4} mb={2}>
            <Text fontSize="xl" fontWeight="bold">{product.price} ₽</Text>
            <Button onClick={handleAddToCart} variant="outline" bg="white" isLoading={adding}>
              Добавить в корзину
            </Button>
          </Flex>
          <Text fontWeight="bold" mb={2}>Описание</Text>
          <Text color="gray.600">{product.description}</Text>
          <Link to="/" style={{ marginTop: 16, display: 'inline-block', color: '#3182ce' }}>← Назад в каталог</Link>
        </Box>
      </Flex>
    </Box>
  );
}
