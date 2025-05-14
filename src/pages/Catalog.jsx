import { Box, Heading, Text, Flex } from '@chakra-ui/react';
import { FaRegLightbulb } from 'react-icons/fa';
import ProductList from '../components/ProductList';
import { products } from '../data';

export default function Catalog() {
  // Заглушка для добавления в корзину
  const handleAddToCart = (product) => {
    alert(`Добавлен в корзину: ${product.name}`);
  };

  return (
    <Box id="catalog" w="100%" minH="100vh">
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
      <ProductList products={products} onAddToCart={handleAddToCart} />
    </Box>
  );
}
