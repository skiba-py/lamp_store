import { Box, Heading, Text, useToast, Spinner } from '@chakra-ui/react';
import { FaRegLightbulb } from 'react-icons/fa';
import ProductList from '../components/ProductList';
import { ordersApi, productsApi } from '../services/api';
import { useState, useEffect } from 'react';

export default function Catalog() {
  const toast = useToast();
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState([]);
  const [fetching, setFetching] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProducts = async () => {
      setFetching(true);
      try {
        const data = await productsApi.getProducts();
        setProducts(data);
      } catch (err) {
        setError('Ошибка при загрузке товаров');
      } finally {
        setFetching(false);
      }
    };
    fetchProducts();
  }, []);

  const handleAddToCart = async (product) => {
    setLoading(true);
    try {
      // Проверяем доступность товара
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
      setLoading(false);
    }
  };

  return (
    <Box id="catalog" w="100%" minH="100vh">
      {(loading || fetching) && (
        <Box position="fixed" top="50%" left="50%" transform="translate(-50%, -50%)" zIndex="1000">
          <Spinner size="xl" color="blue.500" />
        </Box>
      )}
      <Box mb={8} p={6} bg="white" borderRadius="xl" boxShadow="md" display="flex" alignItems="center" gap={4}>
        <Box fontSize="3xl" color="yellow.400"><FaRegLightbulb /></Box>
        <Box>
          <Text fontSize="xl" fontWeight="bold" mb={2}>Добро пожаловать в магазин ламп LAMPS!</Text>
          <Text fontSize="md" color="gray.700">
            В нашем магазине вы найдёте огромный выбор ламп для любого интерьера и настроения: от классических до умных RGB-ламп с управлением через приложение. Мы верим, что свет способен создавать атмосферу уюта, вдохновлять на творчество и даже экономить ваши расходы на электроэнергию!<br/><br/>
            <b>Почему выбирают нас?</b> <br/>
            <ul style={{marginLeft: 20}}>
              <li>💡 Только проверенные бренды и современные технологии</li>
              <li>🌈 Лампы для дома, офиса, растений и вечеринок</li>
              <li>🚚 Быстрая доставка и приятные цены</li>
              <li>🎁 Подарочные упаковки и акции каждую неделю</li>
            </ul>
            <br/>
            Откройте для себя мир света вместе с нами — выберите свою идеальную лампу прямо сейчас!
          </Text>
        </Box>
      </Box>
      <Heading as="h1" size="lg" mb={6} w="100%" textAlign="left">Список товаров</Heading>
      {error && <Text color="red.500" mb={4}>{error}</Text>}
      <ProductList products={products} onAddToCart={handleAddToCart} />
    </Box>
  );
}
