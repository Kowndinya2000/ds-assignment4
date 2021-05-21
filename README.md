# Laptop Recommender

There are 600 laptop models on flipkart.com
````
Please view the screenshots (.png files) to see the execution 
````
main.go performs dual page navigation (2 layer search)
1)  ```- >``` First the code visits flipkart.com and gets all links to laptops in the current page 
2) ```- >``` It indexes to next page and
3)  ```- >``` For each page, it visits all the laptop descripton pages
# Results
Recommended laptops based on users preference are stored in  ```recommendation.json```

Total list of laptops are stored in  ```results.json```
