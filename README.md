# cfn-converter

Archived! Use AWS official tool -> [cfn-flip](https://github.com/awslabs/aws-cfn-template-flip)

Tiny convert tool for AWS CloudFormation template.

Support these functions:

* Convert !Join and !Ref to !Sub
  * Before
    ```yaml
    Resources:
      MyBucket:
        Type: AWS::S3::Bucket
        Properties:
          BucketName: !Join
            - "-"
            - - !Ref AWS::Region
              - !Ref BucketName
    ```
  * After
    ```yaml
    Resources:
      MyBucket:
        Type: AWS::S3::Bucket
        Properties:
          BucketName: !Sub ${AWS::Region}-${BucketName}
    ```

## Download

You can download executable binary from [Release page](https://github.com/taigacat/cfn-converter/releases).

| OS      | binary in Assets                    |
|---------|-------------------------------------|
| Windows | cfn-converter*_windows_amd64.tar.gz |
| MacOS   | cfn-converter*_darwin_amd64.tar.gz  |
| Linux   | cfn-converter*_linux_amd64.tar.gz   |

## Usage

| OS          | command                                               |
|-------------|-------------------------------------------------------|
| Windows     | ```cfn-converter.exe --src <source> --out <output>``` |
| MacOS/Linux | ```cfn-converter --src <source> --out <output>```     |

## Options

| flag     | description                    | default                |
|----------|--------------------------------|------------------------|
| src      | source file                    | -                      |
| out      | output file                    | \<src\>.converted.yaml |
| join2sub | convert !Join and !Ref to !Sub | true                   |
| indent   | indent size                    | 2                      |

